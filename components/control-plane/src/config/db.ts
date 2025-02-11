/*
 * Copyright 2025 Simon Emms <simon@simonemms.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import { registerAs } from '@nestjs/config';
import { TypeOrmModuleOptions } from '@nestjs/typeorm';
import { join } from 'node:path';
import { LoggerOptions } from 'typeorm/logger/LoggerOptions';

export default registerAs('db', (): TypeOrmModuleOptions => {
  /* SSL configuration should default to undefined */
  let ssl: any;
  if (process.env.DB_USE_SSL === 'true') {
    ssl = {
      ca: process.env.DB_SSL_CA,
      cert: process.env.DB_SSL_CERT,
      key: process.env.DB_SSL_KEY,
    };
  }

  let logging: boolean | string | undefined = false;
  const loggingVar = process.env.DB_LOGGING;

  if (loggingVar === 'true') {
    logging = true;
  } else if (loggingVar === 'false') {
    logging = false;
  } else {
    logging = loggingVar;
  }

  return {
    ssl,
    type: (process.env.DB_TYPE as 'mysql' | 'mariadb') ?? 'mysql',
    host: process.env.DB_HOST,
    username: process.env.DB_USERNAME,
    password: process.env.DB_PASSWORD,
    database: process.env.DB_NAME,
    port: process.env.DB_PORT ? Number(process.env.DB_PORT) : undefined,
    migrationsRun: process.env.DB_MIGRATIONS_RUN !== 'false',
    synchronize: process.env.DB_SYNC === 'true',
    autoLoadEntities: true,
    logging: (logging as LoggerOptions) ?? true,
    migrations: [join(__dirname, '..', 'migrations', '*{.ts,.js}')],
  };
});
