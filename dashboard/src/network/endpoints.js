
// Users
const getUsersMe = '/api/users/me'

// Templates
const getTemplate = (id) =>`/api/templates/${id}`
const getTemplates = '/api/templates'
const postTemplates = '/api/templates'
const putTemplates = (id) => `/api/templates/${id}`
const deleleteTemplates = (id) => `/api/templates/${id}`

// Campaigns
const getCampaign = (id) => `/api/campaigns/${id}`
const getCampaignBounces = (id) => `/api/campaigns/${id}/bounces`
const getCampaignClicks = (id) => `/api/campaigns/${id}/clicks`
const getCampaignComplaints = (id) => `/api/campaigns/${id}/complaints`
const getCampaignOpens = (id) => `/api/campaigns/${id}/opens`
const getCampaigns = '/api/campaigns'
const postCampaigns = '/api/campaigns'
const putCampaigns = (id) => `/api/campaigns/${id}`
const deleteCampaigns = (id) => `/api/campaigns/${id}`
const getCampaignsStats = (id) => `/api/campaigns/${id}/stats`
const postCampaignsStart = (id) => `/api/campaigns/${id}/start`

// Groups
const getGroup = (id) => `/api/segments/${id}`
const getGroupSubscribers = (id) => `/api/segments/${id}/subscribers`
const getGroups = '/api/segments'
const postGroups = '/api/segments'
const putGroups = (id) => `/api/segments/${id}`
const deleteGroups = (id) => `/api/segments/${id}`

// Subscribers
const getSubscribers = '/api/subscribers'
const postSubscribers = '/api/subscribers'
const putSubscribers = (id) => `/api/subscribers/${id}`
const deleteSubscribers = (id) => `/api/subscribers/${id}`

// Auth
const signup = '/api/signup'
const signInWithGoogle = 'api/auth/google'
const signInWithGithub = 'api/auth/github'
const signInWithTwitter = 'api/auth/twitter'
const logout = '/api/logout'
const forgotPassword = '/api/forgot-password'
const authenticate = '/api/forgot-password'
const verifyEmail = '/api/verify-email'


export const endpoints = {
    verifyEmail,
    getUsersMe,
    getTemplate,
getTemplates,
postTemplates,
putTemplates,
deleleteTemplates,
getCampaign,
getCampaignBounces,
getCampaignClicks,
getCampaignComplaints,
getCampaignOpens,
getCampaigns,
postCampaigns,
putCampaigns,
deleteCampaigns,
getCampaignsStats,
postCampaignsStart,
getGroup,
getGroupSubscribers,
getGroups,
postGroups,
putGroups,
deleteGroups,
getSubscribers,
postSubscribers,
putSubscribers,
deleteSubscribers,
signup,
signInWithGoogle,
signInWithGithub,
signInWithTwitter,
logout,
forgotPassword,
authenticate
}