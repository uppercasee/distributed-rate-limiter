import { Request, Response, NextFunction } from "express";
import { rateLimiterClient } from "../client";

export function rateLimiterMiddleware(
	req: Request,
	res: Response,
	next: NextFunction,
) {
	const clientId = req.headers["x-client-id"] || "anonymous";

	rateLimiterClient.Check(
		{ client_id: clientId },
		(err: any, response: any) => {
			if (err) {
				console.error("Rate limiter error:", err);
				return res.status(500).send("Rate limiter service unavailable");
			}

			if (response.allowed) {
				next();
			} else {
				console.log("gRPC response:", response);
				res.set("Retry-After", response.retry_after.toString());
				res.status(429).json({
					message: "Rate limit exceeded",
					retry_after: response.retry_after,
				});
			}
		},
	);
}
