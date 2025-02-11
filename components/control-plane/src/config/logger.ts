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
import { ConsoleLoggerOptions, LogLevel } from '@nestjs/common';
import { registerAs } from '@nestjs/config';

// Allow for no logs to be emitted
export const LOG_LEVEL_SILENT = 'silent';

function generateLogLevels(): LogLevel[] {
  const currentLevel = process.env.LOGGER_LEVEL ?? 'log';
  let enabled = false;

  const levels: LogLevel[] = [
    'verbose',
    'debug',
    'log',
    'warn',
    'error',
    'fatal',
  ];

  if (
    !levels.includes(currentLevel as LogLevel) &&
    currentLevel !== LOG_LEVEL_SILENT
  ) {
    throw new Error(`Invalid log level provided: ${currentLevel}`);
  }

  // List the levels in order
  return levels.reduce((levels: LogLevel[], level: LogLevel) => {
    // Once a level is enabled, enable for everything higher than it
    if (level === currentLevel) {
      enabled = true;
    }

    // Once we reach the matching level, go for it
    if (enabled) {
      levels.push(level);
    }
    return levels;
  }, []);
}

export default registerAs(
  'logger',
  (): ConsoleLoggerOptions => ({
    logLevels: generateLogLevels(),
    colors: process.env.LOGGER_COLORS_ENABLED === 'true', // Default to false
    json: process.env.LOGGER_JSON_ENABLED !== 'false', // Default to true
  }),
);
