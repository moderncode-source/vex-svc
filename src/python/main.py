from fastapi import FastAPI, Response
from fastapi.responses import PlainTextResponse

app = FastAPI()


@app.get("/healthz")
async def healthz():
    return Response(status_code=200)


@app.get("/", response_class=PlainTextResponse)
async def read_root():
    return "Hello World"
