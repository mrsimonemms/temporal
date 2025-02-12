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
import { InjectRepository } from '@nestjs/typeorm';
import { Client } from '@temporalio/client';
import { ValidationError, validate } from 'class-validator';
import { Repository } from 'typeorm';

import { EntityValidationException } from '../lib/customExceptions';
import { ProductsService } from '../products/products.service';
import { CreateOrderDto } from './dto/create-order.dto';
import { UpdateOrderDto } from './dto/update-order.dto';
import { Order, OrderStatus } from './entities/order.entity';
import { OrderProduct } from './entities/product.entity';

@Injectable()
export class OrdersService {
  protected readonly logger = new Logger(this.constructor.name);

  @InjectRepository(Order)
  private readonly orderRepository: Repository<Order>;

  @Inject(ProductsService)
  private readonly productService: ProductsService;

  @Inject('ORDER_WORKFLOW_HANDLE')
  private readonly orderWorkflowHandler: Client;

  async create(createOrderDto: CreateOrderDto): Promise<Order> {
    this.logger.debug('Create new order');
    const order = new Order();

    const errs: ValidationError[] = [];
    order.status = OrderStatus.PREPARING;
    order.userId = createOrderDto.userId;
    order.products = await Promise.all(
      createOrderDto.products.map(async (i): Promise<OrderProduct> => {
        const product = await this.productService.findOne(i.productId);

        const o = new OrderProduct();
        o.count = i.count;
        if (product) {
          o.product = product;
        } else {
          const e = new ValidationError();
          e.property = 'productId';
          e.constraints = {
            productId: 'unknown productId',
          };
          e.target = i;
          errs.push(e);
        }

        return o;
      }),
    );

    errs.push(...(await validate(order)));

    if (errs.length > 0) {
      throw new EntityValidationException(errs);
    }

    this.logger.debug('Saving new order');
    const record = await this.orderRepository.save(order);

    return record;
  }

  findAll() {
    return this.orderRepository.find();
  }

  findOne(id: string) {
    return this.orderRepository.findOne({
      where: {
        id,
      },
    });
  }

  update(id: string, updateOrderDto: UpdateOrderDto) {
    return `This action updates a #${id} order`;
  }

  remove(id: string) {
    return `This action removes a #${id} order`;
  }
}
