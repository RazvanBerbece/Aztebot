# syntax=docker/dockerfile:1

# Base Image
FROM eclipse-temurin:17-jdk-jammy

WORKDIR /app

# Copy source files
COPY ./cmd/azteradio-service ./

CMD ["java", "-Dnogui=true", "-jar", "JMusicBot-0.3.9.jar"]