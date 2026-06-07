<template>
  <div class="app-container">
    <div class="background-shapes">
      <div class="shape shape-1"></div>
      <div class="shape shape-2"></div>
      <div class="shape shape-3"></div>
    </div>

    <!-- Main App View -->
    <main v-if="currentUser.token" class="main-content">
      <header class="app-header">
        <div class="logo">
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path
              d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"
            ></path>
            <polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline>
            <line x1="12" y1="22.08" x2="12" y2="12"></line>
          </svg>
          <h2>CRM - Complain Management Module</h2>
        </div>
        <div class="header-right">
          <div class="user-info">
            <span
              class="role-badge"
              :class="{ 'role-admin': currentUser.role === 'admin' }"
            >
              {{ currentUser.role?.toUpperCase() }}
            </span>
            <span class="username">{{ currentUser.username }}</span>
          </div>
          <button @click="logout" class="logout-btn" title="Logout">
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
              <polyline points="16 17 21 12 16 7"></polyline>
              <line x1="21" y1="12" x2="9" y2="12"></line>
            </svg>
          </button>
        </div>
      </header>

      <p class="subtitle">Streamline your customer support experience</p>

      <div class="dashboard-grid">
        <section class="glass-panel form-section">
          <div class="section-header">
            <h3>Log New Complaint</h3>
            <span class="icon">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path
                  d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"
                ></path>
                <path
                  d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"
                ></path>
              </svg>
            </span>
          </div>

          <form @submit.prevent="submitComplaint" class="complaint-form">
            <div class="form-group">
              <label>Customer ID (Shared Link)</label>
              <div class="input-wrapper">
                <input
                  type="number"
                  v-model.number="form.customer_id"
                  required
                  placeholder="e.g. 101"
                />
              </div>
            </div>

            <div class="form-group">
              <label>Complaint Title</label>
              <div class="input-wrapper">
                <input
                  type="text"
                  v-model="form.title"
                  required
                  placeholder="Brief summary of the issue"
                />
              </div>
            </div>

            <div class="form-group">
              <label>Description Details</label>
              <div class="input-wrapper">
                <textarea
                  v-model="form.description"
                  required
                  placeholder="Provide full details here..."
                ></textarea>
              </div>
            </div>

            <button type="submit" class="submit-btn" :disabled="isSubmitting">
              <span v-if="!isSubmitting">Submit Complaint</span>
              <span v-else class="loader"></span>
            </button>
          </form>
        </section>

        <section class="glass-panel list-section">
          <div class="section-header">
            <h3>Active System Complaints</h3>
            <button
              @click="fetchComplaints"
              class="refresh-btn"
              title="Refresh"
            >
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                :class="{ spin: isFetching }"
              >
                <polyline points="23 4 23 10 17 10"></polyline>
                <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path>
              </svg>
            </button>
          </div>

          <div class="table-container">
            <table v-if="complaints.length > 0">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Customer ID</th>
                  <th>Title</th>
                  <th>Status</th>
                  <th v-if="currentUser.role === 'admin'">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in complaints" :key="item.id" class="table-row">
                  <td class="id-col">#{{ item.id }}</td>
                  <td>
                    <strong>{{ item.customer_id }}</strong>
                  </td>
                  <td>
                    <div class="title-cell">{{ item.title }}</div>
                    <div class="desc-cell">{{ item.description }}</div>
                  </td>
                  <td>
                    <span
                      class="status-badge"
                      :class="statusClass(item.status)"
                    >
                      {{ item.status }}
                    </span>
                  </td>
                  <td v-if="currentUser.role === 'admin'" class="actions-cell">
                    <button
                      class="action-btn edit-btn"
                      @click="openStatusModal(item)"
                      title="Change Status"
                    >
                      <svg
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                      >
                        <polyline
                          points="12 3 20 7.5 20 16.5 12 21 4 16.5 4 7.5 12 3"
                        ></polyline>
                        <line x1="12" y1="12" x2="20" y2="7.5"></line>
                        <line x1="12" y1="12" x2="12" y2="21"></line>
                        <line x1="12" y1="12" x2="4" y2="7.5"></line>
                      </svg>
                    </button>
                    <button
                      class="action-btn delete-btn"
                      @click="confirmDelete(item)"
                      title="Delete Complaint"
                    >
                      <svg
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                      >
                        <polyline points="3 6 5 6 21 6"></polyline>
                        <path
                          d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
                        ></path>
                        <line x1="10" y1="11" x2="10" y2="17"></line>
                        <line x1="14" y1="11" x2="14" y2="17"></line>
                      </svg>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>

            <div v-else class="empty-state">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="8" x2="12" y2="12"></line>
                <line x1="12" y1="16" x2="12.01" y2="16"></line>
              </svg>
              <p>No complaints found.</p>
            </div>
          </div>
        </section>
      </div>

      <!-- Status Modal -->
      <div
        v-if="showStatusModal"
        class="modal-overlay"
        @click.self="closeStatusModal"
      >
        <div class="modal-content">
          <div class="modal-header">
            <h4>Change Complaint Status</h4>
            <button class="close-btn" @click="closeStatusModal">×</button>
          </div>
          <div class="modal-body">
            <p class="complaint-info">
              Complaint #{{ selectedComplaint?.id }} -
              {{ selectedComplaint?.title }}
            </p>
            <div class="form-group">
              <label>New Status</label>
              <select v-model="newStatus" class="status-select">
                <option value="Pending">Pending</option>
                <option value="In Progress">In Progress</option>
                <option value="Resolved">Resolved</option>
              </select>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn-secondary" @click="closeStatusModal">
              Cancel
            </button>
            <button class="btn-primary" @click="updateStatus">
              Update Status
            </button>
          </div>
        </div>
      </div>

      <!-- Delete Confirmation Modal -->
      <div
        v-if="showDeleteModal"
        class="modal-overlay"
        @click.self="closeDeleteModal"
      >
        <div class="modal-content danger">
          <div class="modal-header">
            <h4>Delete Complaint</h4>
            <button class="close-btn" @click="closeDeleteModal">×</button>
          </div>
          <div class="modal-body">
            <p class="warning-text">
              Are you sure you want to delete this complaint?
            </p>
            <p class="complaint-info">
              Complaint #{{ selectedComplaint?.id }} -
              {{ selectedComplaint?.title }}
            </p>
            <p class="warning-subtext">This action cannot be undone.</p>
          </div>
          <div class="modal-footer">
            <button class="btn-secondary" @click="closeDeleteModal">
              Cancel
            </button>
            <button class="btn-danger" @click="deleteComplaint">Delete</button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
export default {
  data() {
    return {
      // Auth
      currentUser: { username: "", role: "", token: "" },
      // Complaints
      complaints: [],
      form: { customer_id: null, title: "", description: "" },
      apiBase: "http://localhost:8080/api/complaints",
      isFetching: false,
      isSubmitting: false,
      showStatusModal: false,
      showDeleteModal: false,
      selectedComplaint: null,
      newStatus: "Pending",
    };
  },
  mounted() {
    // Check if logged in from User Management
    const savedToken = localStorage.getItem("crm_token");
    const savedUser = localStorage.getItem("crm_username");
    const savedRole = localStorage.getItem("crm_role");
    
    if (savedToken && savedUser && savedRole) {
      this.currentUser = {
        username: savedUser,
        role: savedRole === "1" ? "admin" : "user",
        token: savedToken
      };
      this.fetchComplaints();
    } else {
      // Not logged in or missing data, redirect to user management
      window.location.href = '/usermgmt/dashboard';
    }
  },
  methods: {
    logout() {
      this.currentUser = { username: "", role: "", token: "" };
      this.complaints = [];
      localStorage.removeItem("crm_token");
      localStorage.removeItem("crm_username");
      localStorage.removeItem("crm_role");
      window.location.href = '/usermgmt/dashboard';
    },
    getAuthHeaders() {
      return {
        "Content-Type": "application/json",
        Authorization: `Bearer ${this.currentUser.token}`,
      };
    },
    async fetchComplaints() {
      this.isFetching = true;
      try {
        const res = await fetch(this.apiBase);
        this.complaints = await res.json();
      } catch (err) {
        console.error("Error fetching integration data:", err);
      } finally {
        setTimeout(() => (this.isFetching = false), 500); // For smooth animation
      }
    },
    async submitComplaint() {
      if (this.isSubmitting) return;
      this.isSubmitting = true;
      try {
        const res = await fetch(this.apiBase, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(this.form),
        });
        if (res.ok) {
          this.form = { customer_id: null, title: "", description: "" };
          await this.fetchComplaints();
        }
      } catch (err) {
        console.error("Error posting complaint:", err);
      } finally {
        this.isSubmitting = false;
      }
    },
    statusClass(status) {
      if (status === "Pending") return "status-pending";
      if (status === "In Progress") return "status-progress";
      return "status-resolved";
    },
    openStatusModal(complaint) {
      this.selectedComplaint = complaint;
      this.newStatus = complaint.status;
      this.showStatusModal = true;
    },
    closeStatusModal() {
      this.showStatusModal = false;
      this.selectedComplaint = null;
      this.newStatus = "Pending";
    },
    async updateStatus() {
      if (!this.selectedComplaint) return;
      try {
        const res = await fetch(
          `${this.apiBase}?id=${this.selectedComplaint.id}`,
          {
            method: "PUT",
            headers: this.getAuthHeaders(),
            body: JSON.stringify({ status: this.newStatus }),
          },
        );
        if (res.ok) {
          this.closeStatusModal();
          await this.fetchComplaints();
        } else if (res.status === 401) {
          alert("Unauthorized - Admin access required");
        }
      } catch (err) {
        console.error("Error updating status:", err);
      }
    },
    confirmDelete(complaint) {
      this.selectedComplaint = complaint;
      this.showDeleteModal = true;
    },
    closeDeleteModal() {
      this.showDeleteModal = false;
      this.selectedComplaint = null;
    },
    async deleteComplaint() {
      if (!this.selectedComplaint) return;
      try {
        const res = await fetch(
          `${this.apiBase}?id=${this.selectedComplaint.id}`,
          {
            method: "DELETE",
            headers: this.getAuthHeaders(),
          },
        );
        if (res.ok) {
          this.closeDeleteModal();
          await this.fetchComplaints();
        } else if (res.status === 401) {
          alert("Unauthorized - Admin access required");
        }
      } catch (err) {
        console.error("Error deleting complaint:", err);
      }
    },
  },
};
</script>

<style>
@import url("https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap");

:root {
  --primary: #6366f1;
  --primary-hover: #4f46e5;
  --bg-color: #0f172a;
  --panel-bg: rgba(30, 41, 59, 0.7);
  --panel-border: rgba(255, 255, 255, 0.1);
  --text-main: #f8fafc;
  --text-muted: #94a3b8;
  --success: #10b981;
  --warning: #f59e0b;
  --info: #3b82f6;
}

* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: "Inter", sans-serif;
  background-color: var(--bg-color);
  color: var(--text-main);
  min-height: 100vh;
  overflow-x: hidden;
}

.app-container {
  position: relative;
  min-height: 100vh;
  padding: 2rem;
  display: flex;
  justify-content: center;
}

/* Dynamic Background Elements */
.background-shapes {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: -1;
  pointer-events: none;
}
.shape {
  position: absolute;
  filter: blur(80px);
  opacity: 0.5;
  border-radius: 50%;
  animation: float 20s infinite ease-in-out;
}
.shape-1 {
  background: #3b82f6;
  width: 400px;
  height: 400px;
  top: -100px;
  left: -100px;
}
.shape-2 {
  background: #8b5cf6;
  width: 300px;
  height: 300px;
  bottom: -50px;
  right: -50px;
  animation-delay: -5s;
}
.shape-3 {
  background: #06b6d4;
  width: 250px;
  height: 250px;
  top: 40%;
  left: 50%;
  animation-delay: -10s;
}

@keyframes float {
  0%,
  100% {
    transform: translate(0, 0) scale(1);
  }
  33% {
    transform: translate(30px, -50px) scale(1.1);
  }
  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }
}

.main-content {
  width: 100%;
  max-width: 1200px;
  z-index: 1;
}

.app-header {
  text-align: center;
  margin-bottom: 3rem;
  animation: slideDown 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}
.app-header .logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 0.5rem;
}
.app-header svg {
  width: 32px;
  height: 32px;
  color: var(--primary);
}
.app-header h2 {
  font-size: 2.2rem;
  font-weight: 700;
  letter-spacing: -0.5px;
  background: linear-gradient(to right, #fff, #94a3b8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
.subtitle {
  color: var(--text-muted);
  font-size: 1.1rem;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 2rem;
}

@media (max-width: 900px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}

/* Glassmorphism Panels */
.glass-panel {
  background: var(--panel-bg);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--panel-border);
  border-radius: 16px;
  padding: 1.5rem;
  box-shadow: 0 20px 40px -15px rgba(0, 0, 0, 0.3);
  animation: fadeIn 1s cubic-bezier(0.16, 1, 0.3, 1);
  transition:
    transform 0.3s ease,
    box-shadow 0.3s ease;
}
.glass-panel:hover {
  transform: translateY(-2px);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.4);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--panel-border);
}
.section-header h3 {
  font-size: 1.25rem;
  font-weight: 600;
}
.section-header .icon svg,
.refresh-btn svg {
  width: 20px;
  height: 20px;
  color: var(--text-muted);
}
.refresh-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  border-radius: 8px;
  transition: all 0.2s;
}
.refresh-btn:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
}
.spin {
  animation: spin 1s linear infinite;
}

/* Form Styles */
.form-group {
  margin-bottom: 1.25rem;
}
.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  margin-bottom: 0.5rem;
  color: var(--text-muted);
}
.input-wrapper input,
.input-wrapper textarea {
  width: 100%;
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid var(--panel-border);
  color: #fff;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  font-family: inherit;
  font-size: 0.95rem;
  transition: all 0.2s;
}
.input-wrapper input:focus,
.input-wrapper textarea:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.2);
}
.input-wrapper textarea {
  min-height: 100px;
  resize: vertical;
}

.submit-btn {
  width: 100%;
  background: var(--primary);
  color: white;
  border: none;
  padding: 0.875rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition:
    background 0.2s,
    transform 0.1s;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 48px;
}
.submit-btn:hover:not(:disabled) {
  background: var(--primary-hover);
}
.submit-btn:active:not(:disabled) {
  transform: scale(0.98);
}
.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* Table Styles */
.table-container {
  overflow-x: auto;
}
table {
  width: 100%;
  border-collapse: collapse;
}
th {
  text-align: left;
  padding: 1rem;
  font-size: 0.875rem;
  color: var(--text-muted);
  font-weight: 500;
  border-bottom: 1px solid var(--panel-border);
}
td {
  padding: 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  vertical-align: top;
}
.table-row {
  transition: background 0.2s;
}
.table-row:hover {
  background: rgba(255, 255, 255, 0.03);
}
.id-col {
  color: var(--text-muted);
  font-family: monospace;
}
.title-cell {
  font-weight: 500;
  margin-bottom: 0.25rem;
}
.desc-cell {
  font-size: 0.875rem;
  color: var(--text-muted);
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 600;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}
.status-pending {
  background: rgba(245, 158, 11, 0.1);
  color: var(--warning);
  border: 1px solid rgba(245, 158, 11, 0.2);
}
.status-progress {
  background: rgba(59, 130, 246, 0.1);
  color: var(--info);
  border: 1px solid rgba(59, 130, 246, 0.2);
}
.status-resolved {
  background: rgba(16, 185, 129, 0.1);
  color: var(--success);
  border: 1px solid rgba(16, 185, 129, 0.2);
}

.empty-state {
  text-align: center;
  padding: 3rem 1rem;
  color: var(--text-muted);
}
.empty-state svg {
  width: 48px;
  height: 48px;
  margin-bottom: 1rem;
  opacity: 0.5;
}

/* Animations */
@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.98);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
@keyframes spin {
  100% {
    transform: rotate(360deg);
  }
}

.loader {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: #fff;
  animation: spin 0.8s linear infinite;
}

/* Action Buttons */
.actions-cell {
  display: flex;
  gap: 0.5rem;
  justify-content: center;
}

.action-btn {
  background: none;
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: var(--text-muted);
  width: 36px;
  height: 36px;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  padding: 0;
}

.action-btn svg {
  width: 16px;
  height: 16px;
}

.action-btn:hover {
  border-color: rgba(255, 255, 255, 0.3);
}

.edit-btn:hover {
  background: rgba(99, 102, 241, 0.1);
  color: var(--primary);
  border-color: var(--primary);
}

.delete-btn:hover {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border-color: #ef4444;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

.modal-content {
  background: var(--panel-bg);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--panel-border);
  border-radius: 12px;
  min-width: 400px;
  max-width: 500px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  animation: slideDown 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.modal-content.danger {
  border-color: rgba(239, 68, 68, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid var(--panel-border);
}

.modal-header h4 {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 1.5rem;
  cursor: pointer;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.2s;
}

.close-btn:hover {
  color: var(--text-main);
}

.modal-body {
  padding: 1.5rem;
}

.complaint-info {
  background: rgba(255, 255, 255, 0.05);
  padding: 0.75rem 1rem;
  border-radius: 6px;
  font-size: 0.875rem;
  margin-bottom: 1rem;
  border-left: 2px solid var(--primary);
}

.warning-text {
  color: #ef4444;
  font-weight: 500;
  margin-bottom: 1rem;
}

.warning-subtext {
  font-size: 0.875rem;
  color: var(--text-muted);
  margin-top: 1rem;
}

.modal-body .form-group {
  margin-bottom: 0;
}

.status-select {
  width: 100%;
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid var(--panel-border);
  color: #fff;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  font-family: inherit;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.2s;
}

.status-select:hover,
.status-select:focus {
  border-color: var(--primary);
  outline: none;
}

.modal-footer {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  padding: 1.5rem;
  border-top: 1px solid var(--panel-border);
}

.btn-primary,
.btn-secondary,
.btn-danger {
  padding: 0.625rem 1.25rem;
  border: none;
  border-radius: 6px;
  font-weight: 600;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: var(--primary);
  color: white;
}

.btn-primary:hover {
  background: var(--primary-hover);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-main);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.3);
}

.btn-danger {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.btn-danger:hover {
  background: rgba(239, 68, 68, 0.3);
  border-color: rgba(239, 68, 68, 0.5);
}

/* Login View */
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  z-index: 1;
  position: relative;
}

.login-panel {
  width: 100%;
  max-width: 420px;
  padding: 2.5rem;
  animation: slideDown 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;
}

.login-header svg {
  width: 48px;
  height: 48px;
  color: var(--primary);
  margin-bottom: 1rem;
}

.login-header h2 {
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.5px;
  margin-bottom: 0.5rem;
  background: linear-gradient(to right, #fff, #94a3b8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.login-header .subtitle {
  color: var(--text-muted);
  font-size: 0.95rem;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.demo-creds {
  margin-top: 1.5rem;
  padding: 1rem;
  background: rgba(99, 102, 241, 0.1);
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: 8px;
  text-align: left;
}

.demo-title {
  font-size: 0.75rem;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 0.75rem;
}

.cred-box {
  font-size: 0.875rem;
  color: var(--text-main);
  line-height: 1.6;
  font-family: monospace;
}

.cred-box strong {
  color: var(--primary);
}

/* Header Modifications */
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  animation: slideDown 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}

.app-header .logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-header svg {
  width: 32px;
  height: 32px;
  color: var(--primary);
}

.app-header h2 {
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.5px;
  background: linear-gradient(to right, #fff, #94a3b8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem 1rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--panel-border);
  border-radius: 8px;
}

.role-badge {
  padding: 0.25rem 0.75rem;
  background: rgba(59, 130, 246, 0.2);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 6px;
  font-size: 0.65rem;
  font-weight: 700;
  letter-spacing: 0.5px;
  color: var(--info);
}

.role-badge.role-admin {
  background: rgba(99, 102, 241, 0.2);
  border-color: rgba(99, 102, 241, 0.3);
  color: var(--primary);
}

.username {
  font-size: 0.875rem;
  color: var(--text-main);
  font-weight: 500;
}

.logout-btn {
  background: none;
  border: 1px solid var(--panel-border);
  color: var(--text-muted);
  width: 40px;
  height: 40px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  padding: 0;
}

.logout-btn svg {
  width: 18px;
  height: 18px;
}

.logout-btn:hover {
  border-color: #ef4444;
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}
</style>
