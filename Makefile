# Copyright 2025 Simon Emms <simon@simonemms.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

COMPONENTS = ./components

cruft-update:
ifeq (,$(wildcard .cruft.json))
	@echo "Cruft not configured"
else
	@cruft check || cruft update --skip-apply-ask --refresh-private-variables
endif
.PHONY: cruft-update

dev:
	$(MAKE) install

	docker compose up --watch
.PHONY: dev

destroy:
	docker compose down
.PHONY: destroy

generate-db-migrations:
	$(shell if [ -z "${NAME}" ]; then echo "NAME must be set"; exit 1; fi)
	docker compose run --rm control-plane npm run migration:generate -- ./src/migrations/${NAME}
.PHONY: generate-db-migrations

install: install-js-deps

install-js-deps:
	@for dir in $(shell ls ${COMPONENTS}/*/package.json); do \
		cd $$(dirname $$dir); \
		echo "Installing $$PWD"; \
		npm ci; \
		cd - > /dev/null; \
	done

	@echo "Installing ${PWD}"
	@npm ci
.PHONY: install-js-deps
