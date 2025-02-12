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
import { ArrayMinSize, IsNotEmpty, Min } from 'class-validator';

export class ProductDTO {
  @ApiProperty({
    description: 'Product ID',
    example: 'c92f19cb-23c7-45bc-ae32-568ee0e33f61',
    required: true,
  })
  @IsNotEmpty()
  productId: string;

  @ApiProperty({
    description: 'Number of items',
    example: 2,
    required: true,
  })
  @IsNotEmpty()
  @Min(1)
  count: number;
}

export class CreateOrderDto {
  @ApiProperty({
    description: 'Products and quantities in the order',
    required: true,
    type: ProductDTO,
    isArray: true,
  })
  @ArrayMinSize(1)
  products: ProductDTO[];

  @ApiProperty({
    description: 'User who owns the order',
    required: true,
    example: 1,
  })
  @IsNotEmpty()
  userId: number;
}
