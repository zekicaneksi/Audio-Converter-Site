FROM node

# Set destination for COPY
WORKDIR /app

# Download dependencies
COPY package*.json ./
RUN npm install

# Copy the files
COPY . .

# Build
RUN npm run build

EXPOSE 3000

# Run
CMD ["npm", "run", "start"]