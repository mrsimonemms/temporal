<!--
  ~ Copyright 2025 Simon Emms <simon@simonemms.com>
  ~
  ~ Licensed under the Apache License, Version 2.0 (the "License");
  ~ you may not use this file except in compliance with the License.
  ~ You may obtain a copy of the License at
  ~
  ~     http://www.apache.org/licenses/LICENSE-2.0
  ~
  ~ Unless required by applicable law or agreed to in writing, software
  ~ distributed under the License is distributed on an "AS IS" BASIS,
  ~ WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  ~ See the License for the specific language governing permissions and
  ~ limitations under the License.
-->

<script lang="ts">
  const providers: { name: string; id: string }[] = [
    {
      name: 'AWS',
      id: 'aws',
    },
  ];

  const regions: { [provider: string]: {} } = {
    aws: ['eu-west-2', 'us-east-1', 'ca-west-1', 'ap-south-1'],
  };

  let provider: string = $state(providers[0].id);
  let region = $state();
  let subnet = $state('10.0.0.0/24');
  let vmCount = $state(1);
</script>

<h1 class="is-size-1">Build cloud resources</h1>

<p class="is-size-5 my-5">Let's build some cloud resources</p>

<form method="POST">
  <div class="field is-horizontal">
    <div class="field-label is-normal">
      <label class="label" for="provider">Provider</label>
    </div>

    <div class="field-body">
      <div class="field">
        <div class="control">
          <div class="select">
            <select bind:value={provider} id="provider">
              {#each providers as item}
                <option value={item.id}>{item.name}</option>
              {/each}
            </select>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div class="field is-horizontal">
    <div class="field-label is-normal">
      <label class="label" for="region">Region</label>
    </div>

    <div class="field-body">
      <div class="field">
        <div class="control">
          <div class="select">
            <select bind:value={region} id="region">
              {#each Object.entries(regions[provider]) as [, value]}
                <option {value}>{value}</option>
              {/each}
            </select>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div class="field is-horizontal">
    <div class="field-label is-normal">
      <label class="label" for="subnet">Subnet</label>
    </div>

    <div class="field-body">
      <div class="field">
        <div class="control">
          <input
            class="input"
            type="text"
            placeholder="10.0.0.0/16"
            value={subnet}
            id="subnet"
            required
          />
        </div>
      </div>
    </div>
  </div>

  <div class="field is-horizontal">
    <div class="field-label is-normal">
      <label class="label" for="vmCount">Number of VMs</label>
    </div>

    <div class="field-body">
      <div class="field">
        <div class="control">
          <input
            class="input"
            type="number"
            placeholder="1"
            value={vmCount}
            min="1"
            max="10"
            id="vmCount"
            required
          />
        </div>
      </div>
    </div>
  </div>

  <input type="submit" class="button is-primary" value="Submit" />
</form>
