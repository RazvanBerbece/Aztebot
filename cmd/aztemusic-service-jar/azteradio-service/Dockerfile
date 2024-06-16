# syntax=docker/dockerfile:1

# Base Image
FROM eclipse-temurin:17-jdk-jammy

WORKDIR /app

# Install utils on container
RUN apt-get update && apt-get install -y iputils-ping && apt install -y make && apt-get install -y vim

# Copy source files
COPY ./cmd/aztemusic-service-jar/JMusicBot-0.3.9.jar ./
COPY ./cmd/aztemusic-service-jar/azteradio-service ./

CMD ["java", "-Dnogui=true", "-jar", "JMusicBot-0.3.9.jar"]