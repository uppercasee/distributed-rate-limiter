FROM node:lts-alpine

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm i

COPY . .

EXPOSE 5174

RUN npm run build

CMD [ "node", "dist/index.js" ]
