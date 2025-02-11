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
import { Controller, Get, Inject, Param } from '@nestjs/common';
import { ApiOkResponse, ApiParam, ApiTags } from '@nestjs/swagger';

import { Product } from './entities/product.entity';
import { ProductsService } from './products.service';

@Controller('products')
@ApiTags('products')
export class ProductsController {
  @Inject(ProductsService)
  private readonly productsService: ProductsService;

  @Get()
  @ApiOkResponse({
    description:
      'List all products - not paginated in this example, but would in reality',
    type: Product,
    isArray: true,
  })
  findAll() {
    return this.productsService.findAll();
  }

  @Get(':id')
  @ApiOkResponse({
    description: 'List single product',
    type: Product,
  })
  @ApiParam({
    name: 'id',
    example: 'c92f19cb-23c7-45bc-ae32-568ee0e33f61',
  })
  findOne(@Param('id') id: string) {
    return this.productsService.findOne(id);
  }
}
