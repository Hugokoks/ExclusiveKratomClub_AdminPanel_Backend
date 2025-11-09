# 1. BUILD FÁZE
FROM golang:1.24-alpine AS builder

# 2. DEFINICE ARGUMENTŮ
# Dockerfile musí definovat, že tyto argumenty přijme z Railway
ARG GITHUB_TOKEN
ARG GO_PRIVATE

# Nastaví pracovní adresář
WORKDIR /app

# 3. INSTALACE NÁSTROJŮ
RUN apk add --no-cache git openssh-client

# 4. KONFIGURACE PRO GO MODULY (ENV)
# GOPRIVATE dostane hodnotu z ARG GO_PRIVATE (např. github.com/Hugokoks/*)
ENV GOPRIVATE=${GO_PRIVATE}
ENV GONOSUM=${GO_PRIVATE}

# 5. AUTORIZACE Gitu přes TOKEN (git config)
# Použije hodnotu z ARG GITHUB_TOKEN (tvůj PAT)
RUN git config --global url."https://oauth2:${GITHUB_TOKEN}@github.com".insteadOf "https://github.com"
# Přesměrování SSH, které už nepoužíváme, ale může pomoci
RUN git config --global url."ssh://git@github.com".insteadOf "https://github.com"

# 6. BUILD
COPY . .
RUN go mod download

# Sestavení binárky
RUN go build -o exclusivekratomclub_adminpanel_backend .

# 7. FINÁLNÍ FÁZE
FROM alpine:latest
WORKDIR /app

COPY --from=builder /exclusivekratomclub_adminpanel_backend /app/exclusivekratomclub_adminpanel_backend

CMD ["/app/exclusivekratomclub_adminpanel_backend"]