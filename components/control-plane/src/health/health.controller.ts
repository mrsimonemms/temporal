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
import { Controller, Get, Inject, VERSION_NEUTRAL } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import {
  HealthCheck,
  HealthCheckService,
  TypeOrmHealthIndicator,
} from '@nestjs/terminus';

import { TemporalService } from '../temporal/temporal.service';

@Controller({
  path: 'health',
  version: VERSION_NEUTRAL,
})
@ApiTags('health')
export class HealthController {
  @Inject(HealthCheckService)
  private health: HealthCheckService;

  @Inject(TemporalService)
  private temporal: TemporalService;

  @Inject(TypeOrmHealthIndicator)
  private db: TypeOrmHealthIndicator;

  @Get()
  @HealthCheck()
  check() {
    // Allow 1 second before timeout
    const timeout = 1000;

    return this.health.check([
      () =>
        this.db.pingCheck('database', {
          timeout,
        }),
      () => this.temporal.healthcheck(timeout),
    ]);
  }
}
