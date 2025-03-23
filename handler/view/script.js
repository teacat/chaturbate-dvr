function data() {
  return {
    settings: {},
    channels: [],
    currentPage: 1,
    itemsPerPage: 5,
    is_updating_channels: false,
    form_data: {
      username: "",
      resolution: "1080",
      resolution_fallback: "down",
      framerate: "30",
      filename_pattern: "{{.Username}}/{{.Year}}-{{.Month}}-{{.Day}}/{{.Username}}_{{.Hour}}-{{.Minute}}-{{.Second}}{{if .Sequence}}_{{.Sequence}}{{end}}",
      split_filesize: 0,
      split_duration: 0,
      interval: 1,
    },

    // Watch for changes in LogLevel
    watchLogLevel() {
      this.$watch("settings.log_level", async (newVal, oldVal) => {
        if (newVal !== oldVal) {
          await this.updateLogLevel();
        }
      });
    },

    // Compute the channels to display for the current page
    get paginatedChannels() {
      const start = (this.currentPage - 1) * this.itemsPerPage;
      return this.channels.slice(start, start + this.itemsPerPage);
    },

    // Calculate total pages
    get totalPages() {
      return Math.ceil(this.channels.length / this.itemsPerPage);
    },

    // Change page on click
    goToPage(page) {
      if (page >= 1 && page <= this.totalPages) {
        this.currentPage = page;
      }
    },
    // openCreateDialog
    openCreateDialog() {
      document.getElementById("create-dialog").showModal();
    },

    // closeCreateDialog
    closeCreateDialog() {
      document.getElementById("create-dialog").close();
      this.resetCreateDialog();
    },

    // submitCreateDialog
    submitCreateDialog() {
      this.createChannel();
      this.closeCreateDialog();
    },

    // error
    error() {
      alert("Error occurred, please refresh the page if something is wrong.");
    },

    //
    async call(path, body) {
      try {
        var resp = await fetch(`/api/${path}`, {
          body: JSON.stringify(body),
          method: "POST",
        });
        if (resp.status !== 200) {
          this.error();
          return [null, true];
        }
        return [await resp.json(), false];
      } catch {
        this.error();
        return [null, true];
      }
    },

    // getSettings
    async getSettings() {
      var [resp, err] = await this.call("get_settings", {});
      if (!err) {
        this.settings = resp;
        this.resetCreateDialog();
        await this.updateLogLevel();
      }
    },

    // init
    async init() {
      document
        .getElementById('create-dialog')
        .addEventListener('close', () => this.resetCreateDialog());

      await this.getSettings(); // Ensure settings are loaded
      this.watchLogLevel(); // Start watching LogLevel after settings load
      await this.listChannels();
      this.listenUpdate();
    },

    async updateLogLevel() {
      const [_, err] = await this.call('update_log_level', {
        log_level: this.settings.log_level,
      });

      if (err) {
        this.error();
      }
    },

    // resetCreateDialog
    resetCreateDialog() {
      document.getElementById("splitting-accordion").open = false;

      // Ensure settings are loaded before resetting form_data
      this.form_data = {
        username: "",
        resolution: this.settings.resolution?.toString() || "1080",
        resolution_fallback: this.settings.resolution_fallback || "down",
        framerate: this.settings.framerate?.toString() || "30",
        filename_pattern: this.settings.filename_pattern || "{{.Username}}/{{.Year}}-{{.Month}}-{{.Day}}/{{.Username}}_{{.Hour}}-{{.Minute}}-{{.Second}}{{if .Sequence}}_{{.Sequence}}{{end}}",
        split_filesize: this.settings.split_filesize?.toString() || "0",
        split_duration: this.settings.split_duration?.toString() || "30",
        interval: this.settings.interval?.toString() || "1",
      };
    },

    // createChannel
    async createChannel() {
      await this.call("create_channel", {
        username: this.form_data.username,
        resolution: parseInt(this.form_data.resolution),
        resolution_fallback: this.form_data.resolution_fallback,
        framerate: parseInt(this.form_data.framerate),
        filename_pattern: this.form_data.filename_pattern,
        split_filesize: parseInt(this.form_data.split_filesize),
        split_duration: parseInt(this.form_data.split_duration),
        interval: parseInt(this.form_data.interval),
      });
    },

    // deleteChannel
    async deleteChannel(username) {
      if (!confirm(`Are you sure you want to delete the channel "${username}"?`)) {
        return;
      }
      var [_, err] = await this.call("delete_channel", { username });
      if (!err) {
        this.channels = this.channels.filter((ch) => ch.username !== username);
      }
    },

    // pauseChannel
    async pauseChannel(username) {
      await this.call("pause_channel", { username });
    },

    // terminateProgram
    async terminateProgram() {
      if (confirm("Are you sure you want to terminate the program?")) {
        alert("The program is terminated, any error messages are safe to ignore.");
        await this.call("terminate_program", {});
      }
    },

    // resumeChannel
    async resumeChannel(username) {
      await this.call("resume_channel", { username });
    },

    // listChannels
    async listChannels() {
      if (this.is_updating_channels) {
        return;
      }
      var [resp, err] = await this.call("list_channels", {});
      if (!err) {
        this.channels = resp.channels;
        this.currentPage = 1;
        this.channels.forEach((ch) => {
          this.scrollLogs(ch.username);
        });
      }
      this.is_updating_channels = false;
    },

    // listenUpdate
    listenUpdate() {
      var source = new EventSource("/api/listen_update");

      source.onmessage = (event) => {
        var data = JSON.parse(event.data);

        // If the channel is not in the list or is stopped, refresh the list.
        if (!this.channels.some((ch) => ch.username === data.username) || data.is_stopped) {
          this.listChannels();
          return;
        }

        var index = this.channels.findIndex((ch) => ch.username === data.username);

        if (index === -1) {
          return;
        }

        this.channels[index].segment_duration = data.segment_duration;
        this.channels[index].segment_filesize = data.segment_filesize;
        this.channels[index].filename = data.filename;
        this.channels[index].last_streamed_at = data.last_streamed_at;
        this.channels[index].is_online = data.is_online;
        this.channels[index].is_paused = data.is_paused;
        this.channels[index].logs = [...this.channels[index].logs, data.log];

        if (this.channels[index].logs.length > 100) {
          this.channels[index].logs = this.channels[index].logs.slice(-100);
        }

        this.scrollLogs(data.username);
      };

      source.onerror = (err) => {
        source.close();
      };
    },

    downloadLogs(username) {
      var a = window.document.createElement("a");
      a.href = window.URL.createObjectURL(new Blob([this.channels[this.channels.findIndex((ch) => ch.username === username)].logs.join("\n")], { type: "text/plain", oneTimeOnly: true }));
      a.download = `${username}_logs.txt`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
    },

    //
    scrollLogs(username) {
      // Wait for the DOM to update.
      setTimeout(() => {
        var logs_element = document.getElementById(`${username}-logs`);

        if (!logs_element) {
          return;
        }
        logs_element.scrollTop = logs_element.scrollHeight;
      }, 1);
    },
  };
}
