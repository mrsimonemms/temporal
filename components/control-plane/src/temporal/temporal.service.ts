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
import { Inject, Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { HealthIndicatorResult, HealthIndicatorStatus } from '@nestjs/terminus';
import { Connection } from '@temporalio/client';
import { grpc as grpcProto } from '@temporalio/proto';
import { setTimeout } from 'timers/promises';

@Injectable()
export class TemporalService {
  protected readonly logger = new Logger(this.constructor.name);

  @Inject(ConfigService)
  private readonly config: ConfigService;

  @Inject('CONNECTION')
  private connection: Connection;

  async onModuleDestroy() {
    this.logger.debug('Disconnecting from Temporal server');
    await this.connection.close();
  }

  async healthcheck(
    timeout: number = 1000,
    serviceName = 'temporal',
  ): Promise<HealthIndicatorResult> {
    let status: HealthIndicatorStatus = 'down';

    try {
      const healthcheck = await Promise.race([
        this.connection.healthService.check({}),
        setTimeout(timeout).then(() => {
          throw new Error('timeout');
        }),
      ]);

      status =
        healthcheck.status ===
        grpcProto.health.v1.HealthCheckResponse.ServingStatus.SERVING
          ? 'up'
          : 'down';
    } catch (err) {
      this.logger.error('Temporal unhealthy', err);
      status = 'down';
    }

    return {
      [serviceName]: {
        status,
      },
    };
  }
}
