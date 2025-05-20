FROM node:lts-alpine

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm i

COPY . .

EXPOSE 5174

RUN npm run build

RUN mkdir -p dist/proto && cp -r src/proto/* dist/proto/

CMD [ "node", "dist/index.js" ]
