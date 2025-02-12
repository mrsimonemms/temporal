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
import { Logger, Provider } from '@nestjs/common';
import {
  Client,
  Workflow,
  WorkflowExecutionAlreadyStartedError,
  WorkflowHandle,
} from '@temporalio/client';

export const taskQueue = 'new-order';
export const workflowId = 'order-workflow';

export const temporalProviders: Provider[] = [
  {
    provide: 'ORDER_WORKFLOW_HANDLE',
    inject: ['WORKFLOW_CLIENT'],
    async useFactory(client: Client): Promise<WorkflowHandle<Workflow>> {
      const logger = new Logger('OrderWorkflowHandle');
      let handle: WorkflowHandle<Workflow>;

      try {
        handle = await client.workflow.start('orderWorkflow', {
          taskQueue,
          workflowId,
        });
        logger.debug('Started new order workflow');
      } catch (err: unknown) {
        if (err instanceof WorkflowExecutionAlreadyStartedError) {
          logger.debug('Reusing existing order workflow');
          handle = client.workflow.getHandle(workflowId);
        } else {
          logger.error('Error creating workflow', { err });
          throw err;
        }
      }

      return handle;
    },
  },
];
