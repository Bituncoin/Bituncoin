FROM node:20-alpine AS frontend

WORKDIR /app

COPY package*.json ./
COPY btng-api/package*.json ./btng-api/

RUN npm ci

COPY . .

RUN cd btng-api && npm run build

FROM mcr.microsoft.com/dotnet/sdk:9.0-alpine AS dotnet-build

WORKDIR /src
COPY src/ .
RUN dotnet publish -c Release -o /app/dotnet

FROM node:20-alpine

WORKDIR /app

COPY --from=frontend /app ./
COPY --from=dotnet-build /app/dotnet ./dotnet

RUN apk add --no-cache dotnet9-runtime python3

COPY BTNGContract.py ./

EXPOSE 3000
EXPOSE 64799

COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

CMD ["/docker-entrypoint.sh"]
