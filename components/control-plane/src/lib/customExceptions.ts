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
import {
  ArgumentsHost,
  Catch,
  ExceptionFilter,
  HttpException,
  HttpStatus,
  Logger,
} from '@nestjs/common';
import { ValidationError } from 'class-validator';
import { FastifyReply } from 'fastify';

export class EntityValidationException extends HttpException {
  errors: ValidationError[];

  constructor(errs: ValidationError[]) {
    super('Validation error', HttpStatus.BAD_REQUEST);

    this.errors = errs;
  }
}

@Catch(EntityValidationException)
export class EntityValidationFilter implements ExceptionFilter {
  private readonly logger = new Logger(this.constructor.name);

  catch(exception: EntityValidationException, host: ArgumentsHost) {
    const response = host.switchToHttp().getResponse<FastifyReply>();

    const statusCode = exception.getStatus();

    this.logger.debug('validation error', exception.errors);

    response.code(statusCode).send({
      statusCode,
      message: exception.message,
      validationErrors: exception.errors,
    });
  }
}
