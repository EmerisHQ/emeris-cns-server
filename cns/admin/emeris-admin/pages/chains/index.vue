<template>
  <div>
    <title-bar :title-stack="titleStack" />
    <hero-bar>
      <div class="columns">
        <div class="column is-half">Chains</div>
        <div class="column">
          <div class="field has-addons">
            <div class="control">
              <file-picker accept="json" v-model="file" />
            </div>
            <div class="control">
              <a
                @click="submit"
                type="submit"
                class="button is-primary disabled"
              >
                Add New Chain
              </a>
            </div>
          </div>
        </div>
      </div>
    </hero-bar>
    <section class="section is-main-section">
      <card-component title="Chains" class="has-table has-mobile-sort-spaced">
        <chains-table />
      </card-component>
      <hr />
    </section>
  </div>
</template>

<script>
import axios from "~/plugins/axios";
import ChainsTable from "@/components/ChainsTable";
import CardComponent from "@/components/CardComponent";
import TitleBar from "@/components/TitleBar";
import HeroBar from "@/components/HeroBar";
import FilePicker from "@/components/FilePicker.vue";

export default {
  name: "Chains",
  components: {
    HeroBar,
    TitleBar,
    CardComponent,
    ChainsTable,
    FilePicker,
  },
  data() {
    return {
      file: "",
    };
  },
  computed: {
    titleStack() {
      return ["Admin", "Chains"];
    },
  },
  methods: {
    async submit() {
      if (!this.file) {
        return;
      }

      let jsontxt = await this.file.text();

      let json = JSON.parse(jsontxt);

      axios
        .post("/add", json)
        .then((res) => {
          this.$buefy.toast.open({
            message: `Successfully uploaded file. Adding chain ${json.chain_name}`,
          });
        })
        .catch((err) => {
          this.$buefy.toast.open({
            message: `Error: ${err.message}`,
            type: "is-danger",
          });
        });
    },
  },
  head() {
    return {
      title: "Chains",
    };
  },
};
</script>
