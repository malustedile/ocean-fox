FROM oven/bun:1.0.25

# Entra no diretório da app dentro do container
WORKDIR /app

# Copia tudo da pasta local `reserva/` para o container
COPY . .

# Instala as dependências
RUN bun install

# Expõe a porta
EXPOSE 3000

# Roda o app
CMD ["bun", "src/index.ts"]
