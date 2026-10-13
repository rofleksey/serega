"""Exercise the built production image using disposable synthetic data only."""
import http.cookiejar
import json
import re
import secrets
import subprocess
import sys
import time
import urllib.error
import urllib.request
import uuid

image = sys.argv[1]
prefix = "serega-smoke-" + uuid.uuid4().hex[:12]
network, database, app = [prefix + suffix for suffix in ("-net", "-db", "-app")]


def docker(*args, data=None):
    result = subprocess.run(["docker", *args], input=data, capture_output=True, timeout=180)
    if result.returncode:
        raise RuntimeError(f"Docker {args[0]} failed: {result.stderr.decode(errors='replace')}")
    return result.stdout


def wait(check, message, seconds=90):
    deadline = time.monotonic() + seconds
    while time.monotonic() < deadline:
        try:
            value = check()
            if value:
                return value
        except (OSError, RuntimeError, KeyError):
            pass
        time.sleep(0.25)
    raise RuntimeError(message)


try:
    docker("network", "create", network)
    docker("run", "-d", "--name", database, "--network", network, "-e", "POSTGRES_HOST_AUTH_METHOD=trust", "-e", "POSTGRES_DB=serega", "postgres:17.6-bookworm")
    wait(lambda: docker("exec", database, "pg_isready", "-h", "127.0.0.1", "-U", "postgres"), "PostgreSQL did not start")
    env = ["-e", f"DATABASE_URL=postgres://postgres@{database}:5432/serega?sslmode=disable", "-e", "SEREGA_COOKIE_SECURE=false", "-e", "OTEL_METRICS_EXPORTER=none"]
    common = ["--network", network, "--read-only", "--cap-drop", "ALL", "--security-opt", "no-new-privileges:true", *env]
    docker("run", "--rm", *common, image, "migrate", "up")
    password = secrets.token_urlsafe(24)
    docker("run", "--rm", "-i", *common, image, "user", "create", "--username", "smoke", "--password-stdin", data=(password + "\n").encode())
    docker("run", "-d", "--name", app, *common, "-p", "127.0.0.1::8080", image)
    port = docker("port", app, "8080").decode().strip().rsplit(":", 1)[1]
    base = "http://127.0.0.1:" + port
    jar = http.cookiejar.CookieJar()
    opener = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(jar))

    def request(path, body=None, extra=None, method=None):
        headers = {"Origin": base}
        if body is not None:
            headers["Content-Type"] = "application/json"
        headers.update(extra or {})
        return opener.open(urllib.request.Request(base + path, data=body, headers=headers, method=method), timeout=10)

    wait(lambda: request("/readyz").status == 204, "Image did not become ready")
    html = request("/")
    assert html.headers["Cache-Control"] == "no-cache"
    asset = re.search(r'src="([^"]+\.js)"', html.read().decode()).group(1)
    assert "immutable" in request(asset).headers["Cache-Control"]
    try:
        request("/assets/missing-chunk.js")
        raise AssertionError("Missing chunk unexpectedly served")
    except urllib.error.HTTPError as error:
        assert error.code == 404
    csrf = json.load(request("/api/v1/auth/csrf"))["token"]
    request("/api/v1/auth/login", json.dumps({"username": "smoke", "password": password}).encode(), extra={"X-CSRF-Token": csrf}).close()
    csrf = json.load(request("/api/v1/auth/csrf"))["token"]
    created = json.load(request("/api/v1/cards", json.dumps({"title": "Smoke card", "description": "Created inside the disposable production image."}).encode(), extra={"X-CSRF-Token": csrf}))
    assert created["status"] == "todo" and created["version"] == 1 and created["createdBy"]["username"] == "smoke"
    moved = json.load(request("/api/v1/cards/" + created["id"], json.dumps({"title": created["title"], "description": created["description"], "status": "doing", "version": created["version"]}).encode(), extra={"X-CSRF-Token": csrf}, method="PATCH"))
    assert moved["status"] == "doing" and moved["version"] == 2
    try:
        request("/api/v1/cards/" + created["id"], json.dumps({"title": created["title"], "description": created["description"], "status": "done", "version": 1}).encode(), extra={"X-CSRF-Token": csrf}, method="PATCH")
        raise AssertionError("Stale update unexpectedly accepted")
    except urllib.error.HTTPError as error:
        assert error.code == 409
    listed = json.load(request("/api/v1/cards"))
    assert any(card["id"] == created["id"] and card["status"] == "doing" for card in listed["cards"])
    print("Production image smoke passed: explicit migration, non-root read-only runtime, login, card create/version conflict/list, asset caching and JSON 404.")
except Exception:
    logs = subprocess.run(["docker", "logs", "--tail", "40", app], capture_output=True, timeout=10)
    if logs.returncode == 0:
        print(logs.stdout.decode(errors="replace") + logs.stderr.decode(errors="replace"), file=sys.stderr)
    raise
finally:
    for name in (app, database):
        subprocess.run(["docker", "rm", "-f", name], capture_output=True, timeout=30)
    subprocess.run(["docker", "network", "rm", network], capture_output=True, timeout=30)
