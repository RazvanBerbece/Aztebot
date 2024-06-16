FROM node:18.16.0-alpine3.17

RUN mkdir -p /app
WORKDIR /app

COPY cmd/azteradio-service/package.json cmd/azteradio-service/package-lock.json .
RUN npm install

COPY cmd/azteradio-service .

CMD [ "node", "src/bot.js"]