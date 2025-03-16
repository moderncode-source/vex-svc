FROM golang:1.23.4-bookworm AS build_image

WORKDIR /app

COPY go.mod ./

RUN go mod download && go mod verify

COPY . .

RUN go generate

# Build a minimal Go binary without debug info,
# with dependencies rebuilt without C libs, rely on Go's only,
# inline function calls aggressively. For options, see:
#   $ go tool compile --help
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux GOAMD64=v2 go build -trimpath -tags osusergo,netgo -o server -a -ldflags="-s -w -buildid=" -gcflags="all=-m=0 -l=2 -dwarf=false" -installsuffix cgo

FROM scratch

COPY --from=build_image /app/server /server

ENTRYPOINT ["/server"]
