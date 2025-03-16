FROM python:3.12.7-slim-bookworm AS build_image

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

WORKDIR /app

RUN pip install --upgrade pip
COPY ./requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir -r /app/requirements.txt

FROM python:3.12.7-slim-bookworm

ENV PYTHONUNBUFFERED=1

WORKDIR /app

COPY --from=build_image /usr/local/lib/python3.12/site-packages /usr/local/lib/python3.12/site-packages
COPY --from=build_image /usr/local/bin /usr/local/bin

COPY . /app

# Note: in a multi-core environment, consider supplying the --workers argument.
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]
