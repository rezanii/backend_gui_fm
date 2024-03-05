FROM devops-registry.smartfren.com/devops/alpine_f:3.12.0

WORKDIR /opt/

RUN mkdir -p logs; \
    mkdir -p config;

COPY backend-gui      /opt/backend-gui

COPY config/.env.staging /opt/config/.env.staging
COPY config/.env.production /opt/config/.env.production

CMD ["/opt/backend-gui", ".env.staging"]
