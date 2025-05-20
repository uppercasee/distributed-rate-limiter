import { Router } from "express";
import { getApiStatus } from "../controllers/apiController";

const router = Router();

router.route("/").get(getApiStatus);

export default router;
