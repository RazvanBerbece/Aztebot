# Base Image
FROM node:18.16.0-alpine3.17

RUN mkdir -p /app
WORKDIR /app

# Setup environment
COPY cmd/azteradio-service .

# Install dependencies
COPY internal/aztemusic-service/package.json internal/aztemusic-service/package-lock.json .
RUN npm install

# Copy rest of source files
COPY internal/aztemusic-service .

# Start the Discord.Js service
CMD [ "node", "src/bot.js"]