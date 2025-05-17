FROM oryd/kratos:v1.3.1


COPY identity.schema.json /etc/config/kratos/identity.schema.json
ENV DSN=postgres://postgres:123456@host.docker.internal:5432/identity?sslmode=disable

#CMD ["-c", "/etc/config/kratos/kratos.yml", "serve"]
#serve command: docker run --rm --name ory-kratos -v D:\workspacego\src\payd\kratos\kratos.yml:/etc/config/kratos/kratos.yml ory-kratos serve -c /etc/config/kratos/kratos.yml
#CMD ["-c", "/etc/config/kratos/kratos.yml", "migrate", "sql"]
#migrate command: docker run --rm --name ory-kratos-migrate -v D:\workspacego\src\payd\kratos\kratos.yml:/etc/config/kratos/kratos.yml ory-kratos migrate sql -e --yes