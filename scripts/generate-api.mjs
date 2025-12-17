import { resolve } from 'path';
import { generateApi } from 'swagger-typescript-api';
import process from 'process';

generateApi({
    name: 'Api.ts',
    output: resolve(process.cwd(), './src/api'),

    input: resolve(process.cwd(), './swagger.yaml'),

    httpClientType: 'axios',
    generateClient: true,
    generateResponses: true,
});