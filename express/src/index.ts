import express from "express";
import api_routes from "./routes/apiRoutes";
import morgan from "morgan";
import cors from "cors";
import helmet from "helmet";
import compression from "compression";

import dotenv from "dotenv";
dotenv.config();

const port: number = 5174;
const app = express();
const API_PREFIX = "/api/v1";

app.use(helmet());
app.use(cors());
app.use(compression());
app.use(morgan("dev"));
app.use(express.json());

// middleware (Rate-Limiter)
// app.use()

app.use(API_PREFIX, api_routes);

const server = app.listen(port, () => {
	console.log(`Server is running http://localhost:${port}`);
});

server.on("error", (error: any) => {
	if (error.code === "EADDRINUSE") {
		console.log(`Port ${port} is already in use. Retrying in 5 seconds...`);
		setTimeout(() => {
			server.close();
			server.listen(port);
		}, 5000);
	} else {
		console.error(error);
	}
});

process.on("SIGINT", () => {
	console.log("\nClosing server...");
	server.close(() => {
		console.log("Server closed.");
		process.exit(0);
	});
});
