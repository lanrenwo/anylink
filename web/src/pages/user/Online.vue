<template>
  <div>
    <el-card>
      <el-form :inline="true">
        <el-form-item>
          <el-select
              v-model="searchCate"
              style="width: 86px;"
              @change="handleSearch">
            <el-option
                label="用户名"
                value="username">
            </el-option>
            <el-option
                label="姓名"
                value="nickname">
            </el-option>
            <el-option
                label="邮箱"
                value="email">
            </el-option>
            <el-option
                label="登录组"
                value="group">
            </el-option>
            <el-option
                label="MAC地址"
                value="mac_addr">
            </el-option>
            <el-option
                label="IP地址"
                value="ip">
            </el-option>
            <el-option
                label="远端地址"
                value="remote_addr">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input
              v-model="searchText"
              placeholder="请输入搜索内容"
              @input="handleSearch">
          </el-input>
        </el-form-item>
        <el-form-item>
          显示休眠用户：
            <el-switch
                v-model="showSleeper"
                @change="handleSearch">
            </el-switch>
        </el-form-item>
        <el-form-item>
          <el-button
              class="extra-small-button"
              type="danger"
              size="mini"
              :loading="loadingOneOffline"
              @click="handleOneOffline">
            一键下线
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-popover
              placement="bottom"
              width="200"
              trigger="click">
            <div>
              <div style="margin-bottom: 10px; font-weight: bold;">选择显示列</div>
              <el-checkbox-group v-model="visibleColumns" @change="saveColumnConfig">
                <el-checkbox
                    v-for="col in allColumns"
                    :key="col.prop"
                    :label="col.prop"
                    style="display: block; margin-bottom: 8px;">
                  {{ col.label }}
                </el-checkbox>
              </el-checkbox-group>
            </div>
            <el-button slot="reference" size="mini" icon="el-icon-setting">
              自定义列
            </el-button>
          </el-popover>
        </el-form-item>
      </el-form>

      <el-table
          ref="multipleTable"
          :data="tableData"
          border>

        <el-table-column
            sortable="true"
            type="index"
            label="序号"
            width="50">
        </el-table-column>

        <el-table-column
            v-if="isColumnVisible('username')"
            prop="username"
            label="用户名"
            min-width="120">
        </el-table-column>

        <el-table-column
            v-if="isColumnVisible('nickname')"
            prop="nickname"
            label="姓名"
            min-width="100">
          <template slot-scope="scope">
            <span v-if="scope.row.nickname">{{ scope.row.nickname }}</span>
            <el-tag v-else type="info" size="mini">未设置</el-tag>
          </template>
        </el-table-column>

        <el-table-column
            v-if="isColumnVisible('email')"
            prop="email"
            label="邮箱"
            min-width="150">
          <template slot-scope="scope">
            <span v-if="scope.row.email">{{ scope.row.email }}</span>
            <el-tag v-else type="info" size="mini">未设置</el-tag>
          </template>
        </el-table-column>

        <el-table-column
            v-if="isColumnVisible('group')"
            prop="group"
            label="登录组"
            min-width="100">
        </el-table-column>

        <el-table-column
            v-if="isColumnVisible('mac_addr')"
            prop="mac_addr"
            label="MAC地址"
            min-width="140">
        </el-table-column>

        <el-table-column
            v-if="isColumnVisible('unique_mac')"
            prop="unique_mac"
            label="唯一MAC"
            width="90">
            <template slot-scope="scope">
                <el-tag v-if="scope.row.unique_mac" type="success">是</el-tag>
                <el-tag v-else type="info">否</el-tag>
            </template>
        </el-table-column>

        <el-table-column
            v-if="isColumnVisible('ip')"
            prop="ip"
            label="IP地址"
            width="140">
        </el-table-column>

        <el-table-column
            v-if="isColumnVisible('remote_addr')"
            prop="remote_addr"
            label="远端地址"
            min-width="140">
        </el-table-column>

        <el-table-column
            v-if="isColumnVisible('transport_protocol')"
            prop="transport_protocol"
            label="传输协议"
            width="90">
        </el-table-column>

        <el-table-column
            v-if="isColumnVisible('tun_name')"
            prop="tun_name"
            label="虚拟网卡"
            width="100">
        </el-table-column>

        <el-table-column
            v-if="isColumnVisible('mtu')"
            prop="mtu"
            label="MTU"
            width="80">
        </el-table-column>

        <el-table-column
            v-if="isColumnVisible('client')"
            prop="is_mobile"
            label="客户端"
            width="80">
          <template slot-scope="scope">
            <i v-if="scope.row.client === 'mobile'" class="el-icon-mobile-phone" style="font-size: 20px;color: red"></i>
            <i v-else class="el-icon-s-platform" style="font-size: 20px;color: blue"></i>
          </template>
        </el-table-column>

        <el-table-column
            v-if="isColumnVisible('bandwidth_realtime')"
            label="实时 上行/下行"
            width="220">
          <template slot-scope="scope">
            <el-tag type="success">{{ scope.row.bandwidth_up }}</el-tag>
            <el-tag class="m-left-10">{{ scope.row.bandwidth_down }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column
            v-if="isColumnVisible('bandwidth_total')"
            label="总量 上行/下行"
            width="200">
          <template slot-scope="scope">
            <el-tag effect="dark" type="success">{{ scope.row.bandwidth_up_all }}</el-tag>
            <el-tag class="m-left-10" effect="dark">{{ scope.row.bandwidth_down_all }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column
            v-if="isColumnVisible('last_login')"
            prop="last_login"
            label="登录时间"
            width="160"
            :formatter="tableDateFormat">
        </el-table-column>

        <el-table-column
            label="操作"
            width="150"
            fixed="right">
          <template slot-scope="scope">
            <el-button
                size="mini"
                type="primary"
                v-if="scope.row.remote_addr !== ''"
                @click="handleReline(scope.row)">重连
            </el-button>

            <el-popconfirm
                class="m-left-10"
                @confirm="handleOffline(scope.row)"
                title="确定要下线用户吗？">
              <el-button
                  slot="reference"
                  size="mini"
                  type="danger">下线
              </el-button>
            </el-popconfirm>

          </template>
        </el-table-column>
      </el-table>

    </el-card>
  </div>
</template>

<script>
import axios from "axios";
import { MessageBox } from 'element-ui';

export default {
  name: "Online",
  components: {},
  mixins: [],
  created() {
    this.$emit('update:route_path', this.$route.path)
    this.$emit('update:route_name', ['用户信息', '在线用户'])
    // 加载保存的列配置
    this.loadColumnConfig()
  },
  mounted() {
    this.getData();

    const chatTimer = setInterval(() => {
      this.getData();
    }, 10000);

    this.$once('hook:beforeDestroy', () => {
      clearInterval(chatTimer);
    })

  },
  data() {
    return {
      tableData: [],
      searchCate: 'username',
      searchText: '',
      showSleeper: false,
      loadingOneOffline: false,
      // 所有可用的列配置
      allColumns: [
        { prop: 'username', label: '用户名' },
        { prop: 'nickname', label: '姓名' },
        { prop: 'email', label: '邮箱' },
        { prop: 'group', label: '登录组' },
        { prop: 'mac_addr', label: 'MAC地址' },
        { prop: 'unique_mac', label: '唯一MAC' },
        { prop: 'ip', label: 'IP地址' },
        { prop: 'remote_addr', label: '远端地址' },
        { prop: 'transport_protocol', label: '传输协议' },
        { prop: 'tun_name', label: '虚拟网卡' },
        { prop: 'mtu', label: 'MTU' },
        { prop: 'client', label: '客户端' },
        { prop: 'bandwidth_realtime', label: '实时带宽' },
        { prop: 'bandwidth_total', label: '总量带宽' },
        { prop: 'last_login', label: '登录时间' },
      ],
      // 默认显示的列
      visibleColumns: [
        'username', 'nickname', 'group', 'ip',
        'remote_addr', 'client', 'bandwidth_realtime', 'last_login'
      ],
    }
  },
  methods: {
    // 判断列是否可见
    isColumnVisible(prop) {
      return this.visibleColumns.includes(prop)
    },
    // 保存列配置到 localStorage
    saveColumnConfig() {
      localStorage.setItem('online_user_columns', JSON.stringify(this.visibleColumns))
    },
    // 加载列配置
    loadColumnConfig() {
      const saved = localStorage.getItem('online_user_columns')
      if (saved) {
        try {
          this.visibleColumns = JSON.parse(saved)
        } catch (e) {
          console.error('加载列配置失败', e)
        }
      }
    },
    handleOffline(row) {
      axios.post('/user/offline?token=' + row.token).then(resp => {
        var data = resp.data
        if (data.code === 0) {
          this.$message.success(data.msg);
          this.getData();
        } else {
          this.$message.error(data.msg);
        }
        console.log(data);
      }).catch(error => {
        this.$message.error('哦，请求出错');
        console.log(error);
      });
    },

    handleReline(row) {
      axios.post('/user/reline?token=' + row.token).then(resp => {
        var data = resp.data
        if (data.code === 0) {
          this.$message.success(data.msg);
          this.getData();
        } else {
          this.$message.error(data.msg);
        }
        console.log(data);
      }).catch(error => {
        this.$message.error('哦，请求出错');
        console.log(error);
      });
    },

    handleEdit(a, row) {
      console.log(a, row)
    },
    handleOneOffline() {
        if (this.tableData === null || this.tableData.length === 0) {
            this.$message.error('错误：当前在线用户表为空，无法执行一键下线操作！');
            return;
        }
        MessageBox.confirm('当前搜索条件下的所有用户将会"下线"，你确定执行吗?', '危险', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'danger'
        }).then(() => {
            try {
                this.loadingOneOffline = true;
                this.getData();
                this.$message.success('操作成功');
                this.loadingOneOffline = false;
                // 清空当前表格
                this.tableData = [];
            } catch (error) {
                this.loadingOneOffline = false;
                this.$message.error('操作失败');
            }
        });
    },
    handleSearch() {
        this.getData();
    },
    getData() {
      axios.get('/user/online',
        {
          params: {
            search_cate: this.searchCate,
            search_text: this.searchText,
            show_sleeper: this.showSleeper,
            one_offline: this.loadingOneOffline
          }
        }
      ).then(resp => {
        var data = resp.data.data
        console.log(data);
        this.tableData = data.datas;
        this.count = data.count
      }).catch(error => {
        this.$message.error('哦，请求出错');
        console.log(error);
      });
    },
  },
}
</script>

<style scoped>
/deep/ .el-form .el-form-item__label,
/deep/ .el-form .el-form-item__content,
/deep/ .el-form .el-input,
/deep/ .el-form .el-select,
/deep/ .el-form .el-button,
/deep/ .el-form .el-select-dropdown__item {
  font-size: 11px;
}
.el-select-dropdown .el-select-dropdown__item {
    font-size: 11px;
    padding: 0 10px;
}
/deep/ .el-input__inner{
    height: 30px;
    padding: 0 10px;
}
.extra-small-button {
  padding: 5px 10px;
}
.m-left-10 {
  margin-left: 10px;
}
</style>
