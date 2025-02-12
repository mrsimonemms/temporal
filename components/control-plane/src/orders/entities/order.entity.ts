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
import { ApiProperty } from '@nestjs/swagger';
import { IsInt, IsNotEmpty } from 'class-validator';
import {
  Column,
  Entity,
  JoinColumn,
  OneToMany,
  PrimaryGeneratedColumn,
} from 'typeorm';

import { TypeormCreateUpdateDeleteTime } from '../../lib/typeorm';
import { OrderProduct } from './product.entity';

export enum OrderStatus {
  PREPARING = 'PREPARING',
}

@Entity()
export class Order extends TypeormCreateUpdateDeleteTime {
  @ApiProperty({
    type: 'string',
    format: 'uuid',
    example: 'c92f19cb-23c7-45bc-ae32-568ee0e33f61',
    required: true,
    description: 'Order ID',
  })
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ApiProperty({
    description: 'Order status',
    enum: OrderStatus,
    default: OrderStatus.PREPARING,
  })
  @Column({
    type: 'enum',
    enum: OrderStatus,
    default: OrderStatus.PREPARING,
  })
  @IsNotEmpty()
  status: OrderStatus;

  // This would link through to a user - out-of-scope
  @ApiProperty({
    description: 'User ID',
    required: true,
    example: 1,
  })
  @Column()
  @IsNotEmpty()
  @IsInt()
  userId: number;

  @OneToMany(() => OrderProduct, (o) => o.order, { eager: true, cascade: true })
  @JoinColumn()
  products: OrderProduct[];
}
