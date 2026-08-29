<template>
  <div id="app">
    <div class="header">
      <h1>区块链学历认证平台</h1>
      <div class="user-info" v-if="user">
        <span>{{ user.role === 'university' ? '高校管理员' : user.role === 'employer' ? '用人单位' : '学生' }}</span>
        <button @click="logout">退出</button>
      </div>
    </div>

    <div class="login" v-if="!user">
      <div class="login-box">
        <h2>登录</h2>
        <input v-model="loginForm.username" placeholder="用户名" />
        <input v-model="loginForm.password" type="password" placeholder="密码" />
        <button @click="login">登录</button>
        <div class="hint">测试账号：admin/admin123 | employer/emp123 | student/stu123</div>
      </div>
    </div>

    <div v-else class="content">
      <!-- 高校管理员：添加学历 -->
      <div v-if="user.role === 'university'" class="panel">
        <h2>添加学历信息</h2>
        <div class="form-row">
          <input v-model="eduForm.certNo" placeholder="证书编号" />
          <input v-model="eduForm.name" placeholder="学生姓名" />
          <input v-model="eduForm.studentId" placeholder="学号" />
        </div>
        <div class="form-row">
          <input v-model="eduForm.school" placeholder="学校" />
          <input v-model="eduForm.major" placeholder="专业" />
          <input v-model="eduForm.degree" placeholder="学位" />
        </div>
        <div class="form-row">
          <input v-model="eduForm.graduationDate" placeholder="毕业日期(YYYY-MM-DD)" />
          <button @click="addEducation">上链存证</button>
        </div>
      </div>

      <!-- 用人单位：验证学历 -->
      <div v-if="user.role === 'employer'" class="panel">
        <h2>验证学历真伪</h2>
        <div class="form-row">
          <input v-model="verifyForm.certNo" placeholder="证书编号" />
          <input v-model="verifyForm.name" placeholder="学生姓名" />
          <button @click="verifyEducation">验证</button>
        </div>
        <div v-if="verifyResult !== null" class="result" :class="verifyResult ? 'success' : 'fail'">
          {{ verifyResult ? '✅ 验证通过：学历信息真实有效' : '❌ 验证失败：未找到匹配记录' }}
        </div>
      </div>

      <!-- 查询学历信息（所有角色） -->
      <div class="panel">
        <h2>查询学历信息</h2>
        <div class="form-row">
          <input v-model="queryForm.certNo" placeholder="输入证书编号查询" />
          <button @click="queryEducation">查询</button>
        </div>
        <div v-if="queryResult" class="result-card">
          <div class="result-item"><span>证书编号：</span>{{ queryResult.certNo }}</div>
          <div class="result-item"><span>姓名：</span>{{ queryResult.name }}</div>
          <div class="result-item"><span>学号：</span>{{ queryResult.studentId }}</div>
          <div class="result-item"><span>学校：</span>{{ queryResult.school }}</div>
          <div class="result-item"><span>专业：</span>{{ queryResult.major }}</div>
          <div class="result-item"><span>学位：</span>{{ queryResult.degree }}</div>
          <div class="result-item"><span>毕业日期：</span>{{ queryResult.graduationDate }}</div>
          <div class="result-item"><span>发证日期：</span>{{ queryResult.issueDate }}</div>
          <div class="result-item"><span>状态：</span>{{ queryResult.status }}</div>
        </div>
      </div>

      <!-- 修改历史 -->
      <div class="panel">
        <h2>修改历史追溯</h2>
        <div class="form-row">
          <input v-model="historyForm.certNo" placeholder="输入证书编号查看历史" />
          <button @click="getHistory">查看历史</button>
        </div>
        <div v-if="historyResult" class="history-list">
          <div v-for="(record, idx) in historyResult" :key="idx" class="history-item">
            <span>交易ID：{{ record.txId.substring(0, 20) }}...</span>
            <span>时间：{{ new Date(record.timestamp).toLocaleString() }}</span>
            <span v-if="record.isDelete">（已删除）</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      user: null,
      token: '',
      loginForm: { username: '', password: '' },
      eduForm: { certNo: '', name: '', studentId: '', school: '', major: '', degree: '', graduationDate: '' },
      verifyForm: { certNo: '', name: '' },
      verifyResult: null,
      queryForm: { certNo: '' },
      queryResult: null,
      historyForm: { certNo: '' },
      historyResult: null,
    }
  },
  methods: {
    async login() {
      const formData = new URLSearchParams()
      formData.append('username', this.loginForm.username)
      formData.append('password', this.loginForm.password)
      const res = await fetch('/api/v1/login', { method: 'POST', body: formData })
      const data = await res.json()
      if (data.code === 0) {
        this.token = data.data.token
        this.user = { role: data.data.role, username: this.loginForm.username }
      } else {
        alert(data.message)
      }
    },
    logout() {
      this.user = null
      this.token = ''
    },
    async addEducation() {
      const res = await fetch('/api/v1/education', {
        method: 'POST',
        headers: { 'Authorization': 'Bearer ' + this.token, 'Content-Type': 'application/json' },
        body: JSON.stringify(this.eduForm)
      })
      const data = await res.json()
      alert(data.message)
    },
    async verifyEducation() {
      const res = await fetch(`/api/v1/education/verify?certNo=${this.verifyForm.certNo}&name=${this.verifyForm.name}`, {
        headers: { 'Authorization': 'Bearer ' + this.token }
      })
      const data = await res.json()
      if (data.code === 0) {
        this.verifyResult = data.data.verified
      } else {
        alert(data.message)
      }
    },
    async queryEducation() {
      const res = await fetch(`/api/v1/education/${this.queryForm.certNo}`, {
        headers: { 'Authorization': 'Bearer ' + this.token }
      })
      const data = await res.json()
      if (data.code === 0) {
        this.queryResult = data.data
      } else {
        alert(data.message)
      }
    },
    async getHistory() {
      const res = await fetch(`/api/v1/education/${this.historyForm.certNo}/history`, {
        headers: { 'Authorization': 'Bearer ' + this.token }
      })
      const data = await res.json()
      if (data.code === 0) {
        this.historyResult = data.data
      } else {
        alert(data.message)
      }
    },
  }
}
</script>

<style>
body { margin: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; }
#app { max-width: 900px; margin: 0 auto; padding: 20px; }
.header { display: flex; justify-content: space-between; align-items: center; padding-bottom: 20px; border-bottom: 2px solid #2563eb; }
.header h1 { font-size: 22px; color: #1a1a2e; }
.login { display: flex; justify-content: center; padding: 60px 0; }
.login-box { background: #fff; padding: 30px; border-radius: 8px; box-shadow: 0 2px 12px rgba(0,0,0,0.1); width: 350px; }
.login-box input { width: 100%; padding: 10px; margin: 8px 0; border: 1px solid #ddd; border-radius: 4px; }
.login-box button { width: 100%; padding: 10px; background: #2563eb; color: white; border: none; border-radius: 4px; cursor: pointer; }
.hint { font-size: 12px; color: #999; margin-top: 10px; }
.content { margin-top: 20px; }
.panel { background: #fff; padding: 20px; border-radius: 8px; margin-bottom: 20px; box-shadow: 0 1px 4px rgba(0,0,0,0.06); }
.panel h2 { font-size: 18px; color: #2563eb; margin-bottom: 15px; }
.form-row { display: flex; gap: 10px; margin-bottom: 10px; flex-wrap: wrap; }
.form-row input { flex: 1; padding: 8px 12px; border: 1px solid #ddd; border-radius: 4px; min-width: 150px; }
.form-row button { padding: 8px 20px; background: #2563eb; color: white; border: none; border-radius: 4px; cursor: pointer; }
.result { padding: 15px; border-radius: 4px; margin-top: 10px; font-weight: 600; }
.result.success { background: #d4edda; color: #155724; }
.result.fail { background: #f8d7da; color: #721c24; }
.result-card { background: #f8fafc; padding: 15px; border-radius: 6px; margin-top: 10px; }
.result-item { padding: 4px 0; }
.result-item span { color: #64748b; font-weight: 600; }
.history-list { margin-top: 10px; }
.history-item { background: #f8fafc; padding: 10px; margin-bottom: 6px; border-radius: 4px; font-size: 13px; display: flex; gap: 20px; }
</style>
