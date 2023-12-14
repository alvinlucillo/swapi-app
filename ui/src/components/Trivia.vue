<!-- The page for searching characters and loading saves searches -->
<template>
  <v-container class="fill-height">
    <v-responsive class="align-center text-center fill-height">
      <v-row>
        <v-col cols="11">
          <v-text-field
            label="Search for a character"
            v-model="searchKey"
            append-icon="mdi-magnify"
            @click:append="searchCharacter"
            :loading="isLoadingSearch"
          >
          </v-text-field>
        </v-col>
        <v-col cols="1">
          <v-btn
            color="primary"
            height="50px"
            @click="saveCurrentSearch"
            :disabled="
              isLoadingSearch ||
              searchResult == undefined ||
              searchResult?.getCharacters?.Characters.length === 0
            "
          >
            Save
          </v-btn>
        </v-col>
      </v-row>
      <v-row style="margin-top: 0; padding-top: 0">
        <v-col cols="11">
          <v-combobox
            clearable
            label="Saved searches"
            :items="savedSearchesItems"
            item-title="SearchKey"
            :item-value="(item) => item.ID"
            density="compact"
            append-icon="mdi-magnify"
            v-model="selectedSavedSearchItem"
            @click:append="loadSavedSearch()"
            :loading="isLoadingLoadSavedSearch"
          ></v-combobox>
        </v-col>
        <v-col cols="1"> </v-col>
      </v-row>

      <!-- Shows the list of characters -->
      <!-- Each character has a table of films and vehicle models -->
      <v-list>
        <v-list-item
          v-for="character in charactersDisplayed"
          :key="character.name"
        >
          <v-card class="border">
            <v-card-title>{{ character.name }}</v-card-title>

            <v-row>
              <v-col cols="6">
                <v-table class="no-border">
                  <thead>
                    <tr>
                      <strong>Films</strong>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="film in character.films" :key="film">
                      <td>{{ film }}</td>
                    </tr>
                    <tr v-if="character.films.length === 0">
                      <td>No films</td>
                    </tr>
                  </tbody>
                </v-table>
              </v-col>

              <v-col cols="6">
                <v-table class="no-border">
                  <thead>
                    <tr>
                      <strong> Vehicle models</strong>
                    </tr>
                  </thead>
                  <tbody>
                    <tr
                      v-for="vehicle in character.vehicleModels"
                      :key="vehicle"
                    >
                      <td>{{ vehicle }}</td>
                    </tr>
                    <tr v-if="character.vehicleModels.length === 0">
                      <td>No vehicles</td>
                    </tr>
                  </tbody>
                </v-table>
              </v-col>
            </v-row>
          </v-card>
        </v-list-item>
      </v-list>
    </v-responsive>

    <!-- Shows snackbar/small dialog messages after each action -->
    <v-snackbar v-model="showSnackbar" :timeout="2000">
      {{ snackBarMessage }}

      <template v-slot:actions>
        <v-btn color="blue" variant="text" @click="showSnackbar = false">
          Close
        </v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<script lang="ts" setup>
import {
  CharactersResult,
  Character,
  SavedSearchesResult,
  SavedSearch,
  useSwapiStore,
} from "@/store/swapi";
import { onMounted } from "vue";
import { ref, computed } from "vue";

const { searchCharacters, getSavedSearches, getSavedSearchByID, saveSearch } =
  useSwapiStore();

// Progress state of text boxes
const isLoadingSearch = ref(false);
const isLoadingLoadSavedSearch = ref(false);

// Snackbar state
const showSnackbar = ref(false);
const snackBarMessage = ref("");

// Value typed in the character search box
const searchKey = ref("");
// Result of character search
const searchResult = ref<CharactersResult>();
// List of saved searches
const savedSearches = ref<SavedSearchesResult>();
// Selected saved search from the dropdown
const selectedSavedSearchItem = ref<SavedSearch>();
// The characters from the search result or saved search
const charactersDisplayed = ref<Character[]>([]);

// Search for a character
const searchCharacter = async () => {
  if (searchKey.value.trim() === "") return;

  selectedSavedSearchItem.value = undefined;

  isLoadingSearch.value = true;
  searchResult.value = await searchCharacters(searchKey.value);
  isLoadingSearch.value = false;

  if (searchResult.value?.getCharacters.Characters.length > 0) {
    charactersDisplayed.value = searchResult.value.getCharacters.Characters;
    showSnackbarMessage("Search completed");
  } else showSnackbarMessage("No results found");
};

// Save the current search
const saveCurrentSearch = async () => {
  if (searchResult.value?.getCharacters.SearchID) {
    await saveSearch(searchResult.value?.getCharacters.SearchID);
    await reloadSavedSearches();

    searchResult.value = undefined;

    showSnackbarMessage("Search saved");
  }
};

// Reload the saved searches after a save
const reloadSavedSearches = async () => {
  savedSearches.value = await getSavedSearches();
};

// Show a snackbar message
const showSnackbarMessage = (message: string) => {
  snackBarMessage.value = message;
  showSnackbar.value = true;
};

// Load a saved search when the dropdown is clicked
const loadSavedSearch = async () => {
  if (selectedSavedSearchItem.value?.ID) {
    searchKey.value = "";
    searchResult.value = undefined;

    isLoadingLoadSavedSearch.value = true;
    const savedSearch = await getSavedSearchByID(
      selectedSavedSearchItem.value?.ID
    );

    if (savedSearch.getSavedSearchesByID?.length > 0) {
      charactersDisplayed.value = savedSearch.getSavedSearchesByID;
      showSnackbarMessage("Search loaded");
    } else showSnackbarMessage("No results found");

    isLoadingLoadSavedSearch.value = false;
  }
};

// Shows the saved searches when the page loads
onMounted(async () => {
  await reloadSavedSearches();
});

// The list of saved searches shown in the dropdown
const savedSearchesItems = computed(
  () => savedSearches.value?.getSavedSearches.map((search) => search) ?? []
);

//
</script>

<style scoped>
.no-border {
  border: none;
}
</style>
