import { Request, Response } from "express";

const getApiStatus = (req: Request, res: Response) => {
	res.send("Api is working fine.");
	console.log("Api is working fine.");
};

export { getApiStatus };
