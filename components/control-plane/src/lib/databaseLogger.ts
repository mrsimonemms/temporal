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
import { Logger as NestLogger } from '@nestjs/common';
import { Logger as TypeOrmLogger } from 'typeorm';

export class DatabaseLogger implements TypeOrmLogger {
  private readonly logger = new NestLogger('SQL');

  logQuery(query: string, parameters?: unknown[]) {
    this.logger.debug(
      {
        query,
        parameters,
      },
      'TypeORM query',
    );
  }

  logQueryError(err: string, query: string, parameters?: unknown[]) {
    this.logger.error(
      {
        err,
        query,
        parameters,
      },
      'TypeORM error',
    );
  }

  logQuerySlow(time: number, query: string, parameters?: unknown[]) {
    this.logger.warn(
      {
        time,
        query,
        parameters,
      },
      'TypeORM slow query',
    );
  }

  logMigration(message: string) {
    this.logger.debug(
      {
        message,
      },
      'TypeORM log migration',
    );
  }

  logSchemaBuild(message: string) {
    this.logger.debug(
      {
        message,
      },
      'TypeORM schema build',
    );
  }

  log(level: 'log' | 'info' | 'warn', message: string) {
    if (level === 'log') {
      return this.logger.log(message);
    }
    if (level === 'info') {
      return this.logger.debug(message);
    }
    if (level === 'warn') {
      return this.logger.warn(message);
    }
  }
}
