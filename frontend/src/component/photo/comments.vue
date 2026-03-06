<template>
  <div class="p-photo-comments">
    <v-card v-if="!readonly" class="mb-4">
      <v-card-text>
        <v-textarea
          v-model="newComment"
          :label="$gettext('Add a comment')"
          rows="3"
          variant="outlined"
          counter="4000"
          maxlength="4000"
        ></v-textarea>
        <v-btn
          color="primary"
          :disabled="!newComment.trim() || loading"
          :loading="loading"
          @click="addComment"
        >
          {{ $gettext('Post Comment') }}
        </v-btn>
      </v-card-text>
    </v-card>

    <v-card v-if="comments.length === 0 && !loading" class="text-center pa-4">
      <v-icon icon="mdi-comment-outline" size="48" color="grey-lighten-2"></v-icon>
      <p class="text-body-1 mt-2">{{ $gettext('No comments yet') }}</p>
      <p class="text-body-2 text-grey">{{ $gettext('Be the first to comment on this photo') }}</p>
    </v-card>

    <v-card v-else-if="comments.length > 0" class="mb-4">
      <v-card-title>
        <v-icon icon="mdi-comment-multiple-outline" class="mr-2"></v-icon>
        {{ $gettext('Comments') }} ({{ comments.length }})
      </v-card-title>
      <v-card-text>
        <v-list>
          <v-list-item
            v-for="comment in comments"
            :key="comment.UID"
            class="mb-2"
          >
            <template v-slot:prepend>
              <v-avatar color="primary" size="40">
                <span class="text-white text-h6">{{ getInitials(comment.UserName) }}</span>
              </v-avatar>
            </template>

            <v-list-item-title class="d-flex align-center">
              <span class="font-weight-bold">{{ comment.UserName || $gettext('Anonymous') }}</span>
              <v-spacer></v-spacer>
              <span class="text-caption text-grey">{{ formatDate(comment.CreatedAt) }}</span>
              <v-btn
                v-if="canDelete(comment)"
                icon="mdi-delete"
                size="small"
                variant="text"
                color="error"
                @click="deleteComment(comment)"
              ></v-btn>
            </v-list-item-title>

            <v-list-item-subtitle class="mt-2">
              <div v-if="comment.Type === 'text'" class="text-body-1">
                {{ comment.Content }}
              </div>
              <div v-else-if="comment.Type === 'audio'">
                <audio controls :src="comment.MediaURL" class="w-100"></audio>
              </div>
              <div v-else-if="comment.Type === 'video'">
                <video controls :src="comment.MediaURL" class="w-100" style="max-height: 200px;"></video>
              </div>
            </v-list-item-subtitle>
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-card>

    <v-snackbar v-model="snackbar.show" :color="snackbar.color">
      {{ snackbar.message }}
    </v-snackbar>
  </div>
</template>

<script>
export default {
  name: "PPhotoComments",
  props: {
    photoUID: {
      type: String,
      required: true,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      comments: [],
      newComment: "",
      loading: false,
      snackbar: {
        show: false,
        message: "",
        color: "success",
      },
      user: this.$session.getUser(),
    };
  },
  mounted() {
    this.loadComments();
  },
  methods: {
    async loadComments() {
      this.loading = true;
      try {
        const response = await this.$api.get(`/photos/${this.photoUID}/comments`);
        this.comments = response.data || [];
      } catch (error) {
        console.error("Failed to load comments:", error);
        this.showError(this.$gettext("Failed to load comments"));
      } finally {
        this.loading = false;
      }
    },
    async addComment() {
      if (!this.newComment.trim()) return;

      this.loading = true;
      try {
        const response = await this.$api.post(`/photos/${this.photoUID}/comments`, {
          content: this.newComment.trim(),
          type: "text",
        });
        this.comments.push(response.data);
        this.newComment = "";
        this.showSuccess(this.$gettext("Comment added"));
      } catch (error) {
        console.error("Failed to add comment:", error);
        this.showError(this.$gettext("Failed to add comment"));
      } finally {
        this.loading = false;
      }
    },
    async deleteComment(comment) {
      if (!confirm(this.$gettext("Are you sure you want to delete this comment?"))) {
        return;
      }

      this.loading = true;
      try {
        await this.$api.delete(`/photos/${this.photoUID}/comments/${comment.UID}`);
        this.comments = this.comments.filter(c => c.UID !== comment.UID);
        this.showSuccess(this.$gettext("Comment deleted"));
      } catch (error) {
        console.error("Failed to delete comment:", error);
        this.showError(this.$gettext("Failed to delete comment"));
      } finally {
        this.loading = false;
      }
    },
    canDelete(comment) {
      if (!this.user) return false;
      return comment.UserUID === this.user.UserUID || this.user.Role === "admin";
    },
    getInitials(name) {
      if (!name) return "?";
      return name
        .split(" ")
        .map(n => n[0])
        .join("")
        .toUpperCase()
        .slice(0, 2);
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString() + " " + date.toLocaleTimeString();
    },
    showSuccess(message) {
      this.snackbar = { show: true, message, color: "success" };
    },
    showError(message) {
      this.snackbar = { show: true, message, color: "error" };
    },
  },
};
</script>

<style scoped>
.p-photo-comments {
  max-width: 800px;
  margin: 0 auto;
}
</style>
