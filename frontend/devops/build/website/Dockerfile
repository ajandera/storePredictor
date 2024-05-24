FROM node:14-alpine3.11 as builder
RUN mkdir /application
COPY . /application
WORKDIR /application
ENV PARTNER_API="https://partner.storepredictor.com/v1/"
ENV INTERNAL_API="https://inter.storepredictor.com/v1/"
ENV CMS_API="https://arualapi.shopycrm.com/v1/public"
ENV VUE_APP_TOKEN="akdfjDF45df/"
ENV VUE_APP_STORAGE="https://arualapi.shopycrm.com"

RUN npm ci && npm run build && npm run generate

FROM node:16-alpine3.11 as application
COPY --from=builder /application /application
WORKDIR /application
CMD ["npm", "run", "start"]
