// byrd-client.js

const helper = {
  // Added userIndex parameter with default value of 0
  setAuthToken: async function (request, userIndex = 0, tokenUrl) {
    try {
      const url = tokenUrl || 'http://localhost:4000/tokens';

      const response = await pm.sendRequest({
        url: url,
        method: 'GET'
      });

      const users = response.json();

      if (users && users.length > userIndex) {
        const token = users[userIndex].value;

        request.headers.remove('Authorization');
        request.headers.add({
          key: 'Authorization',
          value: `Bearer ${token}`
        });

        console.log(`Auth token set successfully using user at index ${userIndex}`);
        return true;
      } else {
        console.error(`No user found at index ${userIndex}`);
        return false;
      }
    } catch (error) {
      console.error('Error setting auth token:', error);
      return false;
    }
  },
  setCompetitorIdFromList: function (response) {
    try {
      const jsonData = response.json();
      if (jsonData.data &&
        jsonData.data.competitors &&
        jsonData.data.competitors.length > 0 &&
        jsonData.data.competitors[0].competitor) {

        const competitorId = jsonData.data.competitors[0].competitor.id;
        pm.collectionVariables.set("competitor_id", competitorId);
        console.log(`Competitor ID set to: ${competitorId}`);
        return true;
      } else {
        console.error('No competitors found in response');
        return false;
      }
    } catch (error) {
      console.error('Error setting competitor ID:', error);
      return false;
    }
  },
  setCompetitorIdFromCreate: function (response) {
    try {
      const jsonData = response.json();
      if (jsonData.data && jsonData.data.id) {
        const competitorId = jsonData.data.id;
        pm.collectionVariables.set("competitor_id", competitorId);
        console.log(`Competitor ID set to: ${competitorId}`);
        return true;
      } else {
        console.error('No competitor ID found in create response');
        return false;
      }
    } catch (error) {
      console.error('Error setting competitor ID:', error);
      return false;
    }
  },
  setWorkspaceIdFromList: function (response) {
    try {
      const jsonData = response.json();
      if (jsonData.data &&
        jsonData.data.workspaces &&
        jsonData.data.workspaces.length > 0) {

        const workspaceId = jsonData.data.workspaces[0].id;
        pm.collectionVariables.set("workspace_id", workspaceId);
        console.log(`Workspace ID set to: ${workspaceId}`);
        return true;
      } else {
        console.error('No workspaces found in response');
        return false;
      }
    } catch (error) {
      console.error('Error setting workspace ID:', error);
      return false;
    }
  },
  setWorkspaceId: function (response) {
    try {
      const jsonData = response.json();
      if (jsonData.data && jsonData.data.id) {
        const workspaceId = jsonData.data.id;
        pm.collectionVariables.set("workspace_id", workspaceId);
        console.log(`Workspace ID set to: ${workspaceId}`);
        return true;
      } else {
        console.error('No workspace ID found in response');
        return false;
      }
    } catch (error) {
      console.error('Error setting workspace ID:', error);
      return false;
    }
  },
  findUserID: function (response, options = {}) {
    try {
      const { workspace_role, membership_status } = options;

      if (!response || !response.data || !response.data.users) {
        console.error('Invalid response structure');
        return null;
      }

      const users = response.data.users;

      // Filter users based on provided criteria
      const matchingUsers = users.filter(user => {
        let matches = true;

        if (workspace_role) {
          matches = matches && user.workspace_role === workspace_role;
        }
        if (membership_status) {
          matches = matches && user.membership_status === membership_status;
        }

        return matches;
      });

      if (matchingUsers.length > 0) {
        const selectedUser = matchingUsers[0];
        console.log(`Found user: ${selectedUser.email} with role: ${selectedUser.workspace_role} and status: ${selectedUser.membership_status}`);
        return selectedUser.user_id;
      } else {
        console.error('No matching user found with specified criteria');
        return null;
      }
    } catch (error) {
      console.error('Error finding user:', error);
      return null;
    }
  },
  findAminUserID: function (response) {
    return findUserID(response, {
      workspace_role: 'admin'
    })
  },
  findPendingUserID: function (response) {
    return findUserID(response, {
      membership_status: 'pending'
    })
  },
  setUserID: function (userID) {
    if (userID) {
      pm.variables.set('user_id', userID);
      console.log('Set userID variable:', userID);
      return true;
    }
    return false;
  },
  setUserIDFromSingleResponse: function (response) {
    try {
      const jsonData = response.json();

      if (jsonData.data && Array.isArray(jsonData.data) && jsonData.data.length > 0) {
        const userId = jsonData.data[0].user_id;

        // Set as both environment and collection variable for flexibility
        pm.environment.set('user_id', userId);
        pm.collectionVariables.set('user_id', userId);

        console.log(`User ID set successfully: ${userId}`);
        return true;
      } else {
        console.error('No user data found in response');
        return false;
      }
    } catch (error) {
      console.error('Error setting user ID:', error);
      return false;
    }
  },
  debugVariableFromRequest: function (variableName) {
    const collectionValue = pm.collectionVariables.get(variableName);
    const envValue = pm.environment.get(variableName);

    console.log(`
Variable: ${variableName}
Collection Value: ${collectionValue || 'not set'}
Environment Value: ${envValue || 'not set'}
Resolved Value: ${collectionValue || envValue || 'not available'}
        `);
  },
  setPageIdFromList: function (response) {
    try {
      const jsonData = response.json();
      if (jsonData.data &&
        jsonData.data.pages &&
        jsonData.data.pages.length > 0) {

        const pageId = jsonData.data.pages[0].id;
        pm.collectionVariables.set("page_id", pageId);
        console.log(`Page ID set from list to: ${pageId}`);
        return true;
      } else {
        console.error('No pages found in list response');
        return false;
      }
    } catch (error) {
      console.error('Error setting page ID from list:', error);
      return false;
    }
  },
  setPageIdFromCreate: function (response) {
    try {
      const jsonData = response.json();
      if (jsonData.data &&
        Array.isArray(jsonData.data) &&
        jsonData.data.length > 0) {

        const pageId = jsonData.data[0].id;
        pm.collectionVariables.set("page_id", pageId);
        console.log(`Page ID set from create to: ${pageId}`);
        return true;
      } else {
        console.error('No page ID found in create response');
        return false;
      }
    } catch (error) {
      console.error('Error setting page ID from create:', error);
      return false;
    }
  },
  setManagementToken: async function (request) {
    try {
      const url = 'http://localhost:4004/token';

      const response = await pm.sendRequest({
        url: url,
        method: 'GET'
      });

      const tokenData = response.json();

      if (tokenData && tokenData.token) {
        request.headers.remove('Authorization');
        request.headers.add({
          key: 'Authorization',
          value: `Bearer ${tokenData.token}`
        });

        console.log('Management token set successfully');
        return true;
      } else {
        console.error('No valid token received from management token server');
        return false;
      }
    } catch (error) {
      console.error('Error setting management token:', error);
      return false;
    }
  },
  // Get the current job ID
  getJobId: function () {
    const jobId = pm.collectionVariables.get("job_id");
    if (!jobId) {
      console.warn('No job ID found in collection variables');
    }
    return jobId;
  },
  setJobIdFromWorkflow: function (response) {
    try {
      const jsonData = response.json();
      if (jsonData.data && jsonData.data.jobID) {
        const jobId = jsonData.data.jobID;
        pm.collectionVariables.set("job_id", jobId);
        console.log(`Job ID set to: ${jobId}`);
        return true;
      } else {
        console.error('No job ID found in workflow response');
        return false;
      }
    } catch (error) {
      console.error('Error setting job ID:', error);
      return false;
    }
  },
  setScheduleIdFromResponse: function (response) {
    try {
      const jsonData = response.json();
      if (jsonData.data && jsonData.data.scheduleID) {
        const scheduleId = jsonData.data.scheduleID;
        pm.collectionVariables.set("schedule_id", scheduleId);
        console.log(`Schedule ID set to: ${scheduleId}`);
        return true;
      } else {
        console.error('No schedule ID found in response');
        return false;
      }
    } catch (error) {
      console.error('Error setting schedule ID:', error);
      return false;
    }
  },
};

module.exports = helper;
