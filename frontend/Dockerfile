FROM node:21.6.2-alpine

WORKDIR /app

COPY package.json ./
COPY yarn.lock ./
COPY . ./
RUN yarn
