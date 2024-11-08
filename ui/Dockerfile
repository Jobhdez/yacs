FROM node:18-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the package.json and package-lock.json files
COPY package*.json ./

# Install the dependencies
RUN npm install

# Copy the rest of the application code into the container
COPY . .


# Install a simple HTTP server to serve static files
RUN npm install -g serve

# Expose the port on which the frontend will run
EXPOSE 5173

# Serve the built application
CMD ["npm", "run", "dev", "--", "--host"]
