FROM node:20-alpine

WORKDIR /app

COPY package*.json ./
COPY btng-api/package*.json ./btng-api/

RUN npm ci

COPY . .

EXPOSE 64799

CMD ["npx", "next", "dev", "btng-api", "-p", "64799", "-H", "0.0.0.0"]
