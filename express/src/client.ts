import * as grpc from "@grpc/grpc-js";
import * as protoLoader from "@grpc/proto-loader";
import path from "path";

const packageDefinition = protoLoader.loadSync(
	path.join(__dirname, "./proto/limiter.proto"),
	{
		keepCase: true,
		longs: String,
		enums: String,
		defaults: true,
		oneofs: true,
	},
);

const protoDescriptor = grpc.loadPackageDefinition(packageDefinition) as any;
const limiterPackage = protoDescriptor.limiter;

export const rateLimiterClient = new limiterPackage.RateLimiterService(
	"grpc-server:50051",
	grpc.credentials.createInsecure(),
);
