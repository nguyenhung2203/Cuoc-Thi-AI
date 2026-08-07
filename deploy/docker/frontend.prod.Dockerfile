FROM node:24.6-alpine AS build
WORKDIR /app
ARG VITE_API_BASE_URL=/api/v1
ARG VITE_WS_URL=
ARG VITE_LIVEKIT_URL
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL VITE_WS_URL=$VITE_WS_URL VITE_LIVEKIT_URL=$VITE_LIVEKIT_URL
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:1.29-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
