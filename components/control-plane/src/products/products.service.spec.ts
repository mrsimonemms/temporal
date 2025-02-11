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
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-unsafe-call */
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
import { Test, TestingModule } from '@nestjs/testing';
import { getRepositoryToken } from '@nestjs/typeorm';
import { Repository } from 'typeorm';

import { Product } from './entities/product.entity';
import { ProductsService } from './products.service';

describe('ProductsService', () => {
  let service: ProductsService;
  let productRepository: any;

  beforeEach(async () => {
    productRepository = {
      find: jest.fn(),
      findOne: jest.fn(),
    };
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        {
          provide: getRepositoryToken(Product),
          useValue: productRepository as Repository<Product>,
        },
        ProductsService,
      ],
    }).compile();

    service = module.get<ProductsService>(ProductsService);
  });

  describe('#findAll', () => {
    it('should return all products', async () => {
      const res = 'some-response';
      productRepository.find.mockResolvedValue(res);

      expect(await service.findAll()).toEqual(res);

      expect(productRepository.find).toHaveBeenCalledWith();
    });
  });

  describe('#findOne', () => {
    it('should return a product by ID', async () => {
      const id = 'some-id';
      const res = 'some-product';
      productRepository.findOne.mockResolvedValue(res);

      expect(await service.findOne(id)).toEqual(res);

      expect(productRepository.findOne).toHaveBeenCalledWith({
        where: {
          id,
        },
      });
    });
  });
});
