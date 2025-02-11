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
import { Response } from 'light-my-request';

import { request } from './setup';

describe('HealthController (e2e)', () => {
  let app: Promise<Response>;

  describe('/', () => {
    describe('GET', () => {
      beforeEach(async () => {
        app = (await request()).inject({
          method: 'GET',
          url: '/health',
        });
      });

      it('should return a healthy state', () =>
        app.then((result) => {
          expect(result.statusCode).toEqual(200);
        }));
    });
  });
});
