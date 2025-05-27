FROM node:23-slim AS builder

WORKDIR /app

COPY . .

RUN npm install
RUN npm run build --mode=production

EXPOSE 3000

CMD ["node", "build"]