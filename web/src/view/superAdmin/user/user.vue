<template>
  <div>
    <warning-bar title="注：右上角头像下拉可切换角色" />
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="用户名">
          <el-input v-model="searchInfo.username" placeholder="用户名" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="searchInfo.nickname" placeholder="昵称" />
        </el-form-item>
        <el-form-item label="账号类型">
          <el-select v-model="searchInfo.user_type" placeholder="请选择账号类型" clearable>
            <el-option label="系统用户" :value="1" />
            <el-option label="商户用户" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="searchInfo.phone" placeholder="手机号" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="searchInfo.email" placeholder="邮箱" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">
            查询
          </el-button>
          <el-button icon="refresh" @click="onReset"> 重置 </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="addUser"
        >新增用户</el-button
        >
      </div>
      <el-table :data="tableData" row-key="ID">
        <el-table-column align="left" label="头像" min-width="75">
          <template #default="scope">
            <CustomPic style="margin-top: 8px" :pic-src="scope.row.headerImg" />
          </template>
        </el-table-column>
        <el-table-column align="left" label="ID" min-width="50" prop="ID" />
        <el-table-column
            align="left"
            label="用户名"
            min-width="150"
            prop="userName"
        />
        <el-table-column
            align="left"
            label="昵称"
            min-width="150"
            prop="nickName"
        />
        <el-table-column
            align="left"
            label="账号类型"
            min-width="120"
        >
          <template #default="scope">
            <el-tag :type="scope.row.userType === 2 ? 'success' : 'primary'">
              {{ scope.row.userType === 2 ? '商户用户' : '系统用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="商户"
            min-width="150"
            prop="mchName"
        />
        <el-table-column
            align="left"
            label="手机号"
            min-width="180"
            prop="phone"
        />
        <el-table-column
            align="left"
            label="邮箱"
            min-width="180"
            prop="email"
        />
        <el-table-column align="left" label="用户角色" min-width="200">
          <template #default="scope">
            <el-cascader
                v-model="scope.row.authorityIds"
                :options="authOptions"
                :show-all-levels="false"
                collapse-tags
                :props="{
                multiple: true,
                checkStrictly: true,
                label: 'authorityName',
                value: 'authorityId',
                disabled: 'disabled',
                emitPath: false
              }"
                :clearable="false"
                @visible-change="
                (flag) => {
                  changeAuthority(scope.row, flag, 0)
                }
              "
                @remove-tag="
                (removeAuth) => {
                  changeAuthority(scope.row, false, removeAuth)
                }
              "
            />
          </template>
        </el-table-column>
        <el-table-column align="left" label="启用" min-width="150">
          <template #default="scope">
            <el-switch
                v-model="scope.row.enable"
                inline-prompt
                :active-value="1"
                :inactive-value="2"
                @change="
                () => {
                  switchEnable(scope.row)
                }
              "
            />
          </template>
        </el-table-column>

        <el-table-column label="操作" :min-width="appStore.operateMinWith" fixed="right">
          <template #default="scope">
            <el-button
                type="primary"
                link
                icon="delete"
                @click="deleteUserFunc(scope.row)"
            >删除</el-button
            >
            <el-button
                type="primary"
                link
                icon="edit"
                @click="openEdit(scope.row)"
            >编辑</el-button
            >
            <el-button
                type="primary"
                link
                icon="magic-stick"
                @click="resetPasswordFunc(scope.row)"
            >重置密码</el-button
            >
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>
    <!-- 重置密码对话框 -->
    <el-dialog
        v-model="resetPwdDialog"
        title="重置密码"
        width="500px"
        :close-on-click-modal="false"
        :close-on-press-escape="false"
    >
      <el-form :model="resetPwdInfo" ref="resetPwdForm" label-width="100px">
        <el-form-item label="用户账号">
          <el-input v-model="resetPwdInfo.userName" disabled />
        </el-form-item>
        <el-form-item label="用户昵称">
          <el-input v-model="resetPwdInfo.nickName" disabled />
        </el-form-item>
        <el-form-item label="账号类型">
          <el-input :value="resetPwdInfo.userType === 2 ? '商户用户' : '系统用户'" disabled />
        </el-form-item>
        <el-form-item label="新密码">
          <div class="flex w-full">
            <el-input class="flex-1" v-model="resetPwdInfo.password" placeholder="请输入新密码" show-password />
            <el-button type="primary" @click="generateRandomPassword" style="margin-left: 10px">
              生成随机密码
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeResetPwdDialog">取 消</el-button>
          <el-button type="primary" @click="confirmResetPassword">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-drawer
        v-model="addUserDialog"
        :size="appStore.drawerSize"
        :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">用户</span>
          <div>
            <el-button @click="closeAddUserDialog">取 消</el-button>
            <el-button type="primary" @click="enterAddUserDialog"
            >确 定</el-button
            >
          </div>
        </div>
      </template>

      <el-form
          ref="userForm"
          :rules="rules"
          :model="userInfo"
          label-width="80px"
      >
        <el-form-item
            v-if="dialogFlag === 'add'"
            label="用户名"
            prop="userName"
        >
          <el-input v-model="userInfo.userName" />
        </el-form-item>
        <el-form-item v-if="dialogFlag === 'add'" label="密码" prop="password">
          <el-input v-model="userInfo.password" show-password />
        </el-form-item>
        <el-form-item label="账号类型" prop="userType">
          <el-radio-group v-model="userInfo.userType" @change="handleUserTypeChange">
            <el-radio :label="1">系统用户</el-radio>
            <el-radio :label="2">商户用户</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="昵称" prop="nickName">
          <el-input v-model="userInfo.nickName" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="userInfo.phone" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userInfo.email" />
        </el-form-item>

        <!-- 商户选择框，仅当userType为2（商户用户）且当前登录用户为系统用户时显示 -->
        <el-form-item
            v-if="userInfo.userType === 2 && showMerchantField"
            label="关联商户"
            prop="mchNo"
            :rules="userInfo.userType === 2 ? [{ required: true, message: '请选择商户', trigger: 'blur' }] : []"
        >
          <el-select
              v-model="userInfo.mchNo"
              filterable
              placeholder="请选择商户"
              clearable
              remote
              style="width: 100%"
              :remote-method="remoteSearchVenuer"
              :loading="venuerLoading"
              @focus="loadMerchantOptions"
          >
            <!-- 如果是在编辑模式且已有商户信息，显示当前商户 -->
            <el-option
                v-if="dialogFlag === 'edit' && userInfo.mchNo && currentMerchantInfo"
                :key="currentMerchantInfo.mch_no"
                :label="currentMerchantInfo.mch_name"
                :value="currentMerchantInfo.mch_no"
                style="color: #409EFF;"
            >
              <span style="float: left">{{ currentMerchantInfo.mch_name }}</span>
              <span style="float: right; color: #67C23A; font-size: 12px">当前商户</span>
            </el-option>

            <!-- 其他商户选项 -->
            <el-option
                v-for="item in venuerList"
                :key="item.mch_no"
                :label="item.mch_name"
                :value="item.mch_no"
            >
              <span style="float: left">{{ item.mch_name }}</span>
              <span style="float: right; color: #8492a6; font-size: 13px">{{ item.mch_no }}</span>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="用户角色" prop="authorityId">
          <el-cascader
              v-model="userInfo.authorityIds"
              style="width: 100%"
              :options="authOptions"
              :show-all-levels="false"
              :props="{
              multiple: true,
              checkStrictly: true,
              label: 'authorityName',
              value: 'authorityId',
              disabled: 'disabled',
              emitPath: false
            }"
              :clearable="false"
          />
        </el-form-item>
        <el-form-item label="启用" prop="disabled">
          <el-switch
              v-model="userInfo.enable"
              inline-prompt
              :active-value="1"
              :inactive-value="2"
          />
        </el-form-item>
        <el-form-item label="头像" label-width="80px">
          <SelectImage v-model="userInfo.headerImg" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  getUserList,
  setUserAuthorities,
  register,
  deleteUser
} from '@/api/user'

import {getAuthorityList} from '@/api/authority'
import CustomPic from '@/components/customPic/index.vue'
import WarningBar from '@/components/warningBar/warningBar.vue'
import {setUserInfo, resetPassword, getUserInfo} from '@/api/user.js'

import {nextTick, ref, watch} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import SelectImage from '@/components/selectImage/selectImage.vue'
import {useAppStore} from "@/pinia";
import {getMerchantList} from "@/plugin/payment/api/merchant.js";

defineOptions({
  name: 'User'
})

const appStore = useAppStore()

// 添加当前用户信息变量
const currentUserInfo = ref({})
const showMerchantField = ref(false)
const currentMerchantInfo = ref(null)

// 获取当前登录用户信息
const fetchCurrentUserInfo = async () => {
  try {
    const res = await getUserInfo()
    if (res.code === 0) {
      currentUserInfo.value = res.data.userInfo
      console.log('当前用户信息:', currentUserInfo.value)

      // 判断是否显示商户字段：只有系统用户(userType !== 2)才能管理商户
      showMerchantField.value = currentUserInfo.value.userType !== 2
    } else {
      console.error('获取用户信息失败:', res.msg)
      showMerchantField.value = false
    }
  } catch (error) {
    console.error('获取用户信息异常:', error)
    showMerchantField.value = false
  }
}

const venuerLoading = ref(false)
const venuerList = ref([])

// 加载商户选项
const loadMerchantOptions = async () => {
  // 只有在商户用户类型且当前登录用户是系统用户时才需要加载
  if (userInfo.value.userType === 2 && showMerchantField.value) {
    try {
      // 清空列表
      venuerList.value = []

      // 如果是编辑模式且已有商户信息，先获取当前商户信息
      if (dialogFlag.value === 'edit' && userInfo.value.mchNo) {
        try {
          const res = await getMerchantList({mch_no: userInfo.value.mchNo, page: 1, page_size: 1})
          if (res.code === 0 && res.data?.list?.length > 0) {
            currentMerchantInfo.value = res.data.list[0]
            // 将当前商户添加到列表
            venuerList.value.push(currentMerchantInfo.value)
          }
        } catch (error) {
          console.error('加载当前商户信息失败:', error)
        }
      }

      // 加载其他商户（如果需要的话）
      if (dialogFlag.value === 'add' || venuerList.value.length === 0) {
        const res = await getMerchantList({page: 1, page_size: 10})
        if (res.code === 0 && res.data?.list) {
          // 过滤掉已经在列表中的商户
          const existingMchNos = new Set(venuerList.value.map(item => item.mch_no))
          const newMerchants = res.data.list.filter(item => !existingMchNos.has(item.mch_no))
          venuerList.value = [...venuerList.value, ...newMerchants]
        }
      }
    } catch (error) {
      console.error('加载商户列表失败:', error)
      ElMessage.error('加载商户列表失败')
    }
  }
}

// 远程搜索商户
const remoteSearchVenuer = async (query) => {
  if (userInfo.value.userType !== 2 || !showMerchantField.value) return

  venuerLoading.value = true
  try {
    const res = await getMerchantList({mch_name: query, page: 1, page_size: 10})
    if (res.code === 0) {
      const newMerchants = res.data.list || []
      // 去重：过滤掉已经在列表中的商户
      const existingMchNos = new Set(venuerList.value.map(item => item.mch_no))
      const filteredMerchants = newMerchants.filter(item => !existingMchNos.has(item.mch_no))

      // 如果是编辑模式且当前商户不在搜索结果中，确保当前商户在列表中
      if (dialogFlag.value === 'edit' && currentMerchantInfo.value) {
        const hasCurrentMerchant = filteredMerchants.some(item => item.mch_no === currentMerchantInfo.value.mch_no)
        if (!hasCurrentMerchant && !venuerList.value.some(item => item.mch_no === currentMerchantInfo.value.mch_no)) {
          venuerList.value.unshift(currentMerchantInfo.value)
        }
      }

      venuerList.value = [...venuerList.value, ...filteredMerchants]
    } else {
      ElMessage.error(res.msg || '搜索商户失败')
    }
  } finally {
    venuerLoading.value = false
  }
}

// 用户类型变更处理
const handleUserTypeChange = (value) => {
  console.log('用户类型变更:', value)
  // 当用户类型从商户用户改为系统用户时，清空商户选择
  if (value !== 2) {
    userInfo.value.mchNo = ''
    currentMerchantInfo.value = null
  } else {
    // 当用户类型改为商户用户时，如果是编辑模式且已有商户信息，加载商户信息
    if (dialogFlag.value === 'edit' && userInfo.value.mchNo) {
      loadMerchantOptions()
    }
  }
}

const searchInfo = ref({
  username: '',
  nickname: '',
  phone: '',
  email: '',
  user_type: ''
})

const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = {
    username: '',
    nickname: '',
    phone: '',
    email: '',
    user_type: ''
  }
  getTableData()
}

// 初始化相关
const setAuthorityOptions = (AuthorityData, optionsData) => {
  AuthorityData &&
  AuthorityData.forEach((item) => {
    if (item.children && item.children.length) {
      const option = {
        authorityId: item.authorityId,
        authorityName: item.authorityName,
        children: []
      }
      setAuthorityOptions(item.children, option.children)
      optionsData.push(option)
    } else {
      const option = {
        authorityId: item.authorityId,
        authorityName: item.authorityName
      }
      optionsData.push(option)
    }
  })
}

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async () => {
  const table = await getUserList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

watch(
    () => tableData.value,
    () => {
      setAuthorityIds()
    }
)

const authOptions = ref([])
const setOptions = (authData) => {
  authOptions.value = []
  setAuthorityOptions(authData, authOptions.value)
}

const initPage = async () => {
  getTableData()
  const res = await getAuthorityList()
  setOptions(res.data)
  // 初始化时获取当前用户信息
  await fetchCurrentUserInfo()
}

initPage()

// 重置密码对话框相关
const resetPwdDialog = ref(false)
const resetPwdForm = ref(null)
const resetPwdInfo = ref({
  ID: '',
  userName: '',
  nickName: '',
  userType: '',
  password: '',
})

// 生成随机密码
const generateRandomPassword = () => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*'
  let password = ''
  for (let i = 0; i < 12; i++) {
    password += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  resetPwdInfo.value.password = password
  // 复制到剪贴板
  navigator.clipboard.writeText(password).then(() => {
    ElMessage({
      type: 'success',
      message: '密码已复制到剪贴板'
    })
  }).catch(() => {
    ElMessage({
      type: 'error',
      message: '复制失败，请手动复制'
    })
  })
}

// 打开重置密码对话框
const resetPasswordFunc = (row) => {
  resetPwdInfo.value.ID = row.ID
  resetPwdInfo.value.userName = row.userName
  resetPwdInfo.value.nickName = row.nickName
  resetPwdInfo.value.userType = row.userType
  resetPwdInfo.value.password = ''
  resetPwdDialog.value = true
}

// 确认重置密码
const confirmResetPassword = async () => {
  if (!resetPwdInfo.value.password) {
    ElMessage({
      type: 'warning',
      message: '请输入或生成密码'
    })
    return
  }

  const res = await resetPassword({
    ID: resetPwdInfo.value.ID,
    password: resetPwdInfo.value.password
  })

  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: res.msg || '密码重置成功'
    })
    resetPwdDialog.value = false
  } else {
    ElMessage({
      type: 'error',
      message: res.msg || '密码重置失败'
    })
  }
}

// 关闭重置密码对话框
const closeResetPwdDialog = () => {
  resetPwdInfo.value.password = ''
  resetPwdDialog.value = false
}

const setAuthorityIds = () => {
  tableData.value &&
  tableData.value.forEach((user) => {
    user.authorityIds =
        user.authorities &&
        user.authorities.map((i) => {
          return i.authorityId
        })
  })
}

const deleteUserFunc = async (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await deleteUser({id: row.ID})
    if (res.code === 0) {
      ElMessage.success('删除成功')
      await getTableData()
    }
  })
}

// 弹窗相关
const userInfo = ref({
  userName: '',
  password: '',
  nickName: '',
  headerImg: '',
  authorityId: '',
  authorityIds: [],
  enable: 1,
  userType: 1, // 默认系统用户
  mchNo: ''
})

const rules = ref({
  userName: [
    {required: true, message: '请输入用户名', trigger: 'blur'},
    {min: 5, message: '最低5位字符', trigger: 'blur'}
  ],
  password: [
    {required: true, message: '请输入用户密码', trigger: 'blur'},
    {min: 6, message: '最低6位字符', trigger: 'blur'}
  ],
  nickName: [{required: true, message: '请输入用户昵称', trigger: 'blur'}],
  userType: [
    {required: true, message: '请选择账号类型', trigger: 'blur'}
  ],
  phone: [
    {
      pattern: /^1([38][0-9]|4[014-9]|[59][0-35-9]|6[2567]|7[0-8])\d{8}$/,
      message: '请输入合法手机号',
      trigger: 'blur'
    }
  ],
  email: [
    {
      pattern: /^([0-9A-Za-z\-_.]+)@([0-9a-z]+\.[a-z]{2,3}(\.[a-z]{2})?)$/g,
      message: '请输入正确的邮箱',
      trigger: 'blur'
    }
  ],
  authorityId: [
    {required: true, message: '请选择用户角色', trigger: 'blur'}
  ]
})

const userForm = ref(null)
const enterAddUserDialog = async () => {
  userInfo.value.authorityId = userInfo.value.authorityIds[0]
  userForm.value.validate(async (valid) => {
    if (valid) {
      const req = {
        ...userInfo.value
      }

      // 如果是系统用户，清空商户信息
      if (req.userType === 1) {
        req.mchNo = ''
      }

      if (dialogFlag.value === 'add') {
        const res = await register(req)
        if (res.code === 0) {
          ElMessage({type: 'success', message: '创建成功'})
          await getTableData()
          closeAddUserDialog()
        }
      }
      if (dialogFlag.value === 'edit') {
        const res = await setUserInfo(req)
        if (res.code === 0) {
          ElMessage({type: 'success', message: '编辑成功'})
          await getTableData()
          closeAddUserDialog()
        }
      }
    }
  })
}

const addUserDialog = ref(false)
const closeAddUserDialog = () => {
  userForm.value.resetFields()
  userInfo.value.headerImg = ''
  userInfo.value.authorityIds = []
  userInfo.value.userType = 1
  userInfo.value.mchNo = ''
  currentMerchantInfo.value = null
  venuerList.value = []
  addUserDialog.value = false
}

const dialogFlag = ref('add')

const addUser = () => {
  dialogFlag.value = 'add'
  userInfo.value = {
    userName: '',
    password: '',
    nickName: '',
    headerImg: '',
    authorityId: '',
    authorityIds: [],
    enable: 1,
    userType: 1,
    mchNo: ''
  }
  currentMerchantInfo.value = null
  venuerList.value = []
  addUserDialog.value = true
}

const tempAuth = {}
const changeAuthority = async (row, flag, removeAuth) => {
  if (flag) {
    if (!removeAuth) {
      tempAuth[row.ID] = [...row.authorityIds]
    }
    return
  }
  await nextTick()
  const res = await setUserAuthorities({
    ID: row.ID,
    authorityIds: row.authorityIds
  })
  if (res.code === 0) {
    ElMessage({type: 'success', message: '角色设置成功'})
  } else {
    if (!removeAuth) {
      row.authorityIds = [...tempAuth[row.ID]]
      delete tempAuth[row.ID]
    } else {
      row.authorityIds = [removeAuth, ...row.authorityIds]
    }
  }
}

const openEdit = async (row) => {
  dialogFlag.value = 'edit'
  userInfo.value = JSON.parse(JSON.stringify(row))

  // 确保userType存在
  if (!userInfo.value.userType) {
    userInfo.value.userType = 1
  }

  // 如果用户有商户ID，设置到表单中
  if (row.mchNo) {
    userInfo.value.mchNo = row.mchNo
    console.log('用户已关联商户ID:', row.mchNo, '用户类型:', userInfo.value.userType)
  }

  addUserDialog.value = true
}

const switchEnable = async (row) => {
  userInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    ...userInfo.value
  }
  const res = await setUserInfo(req)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: `${req.enable === 2 ? '禁用' : '启用'}成功`
    })
    await getTableData()
    userInfo.value.headerImg = ''
    userInfo.value.authorityIds = []
  }
}
</script>

<style lang="scss">
.header-img-box {
  @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
}

/* 商户选项样式优化 */
.el-select-dropdown__item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>