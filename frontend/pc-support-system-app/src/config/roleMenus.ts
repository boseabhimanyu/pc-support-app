export type MenuItem = {
    label: string;
    path?: string;
    children?: MenuItem[];
};

export const roleMenus: Record<string, MenuItem[]> = {

    receptionist: [

         {
            label: "Dashboard",
            path: "/receptionist",
        },

        {
            label: "Customer Management",
            path: "/receptionist/customers",
        },

        {
            label: "Job Management",
            path: "/receptionist/jobs",
        },

        {
            label: "My Profile",
            path: "/receptionist/profile",
        },

    ],

    technician: [
         {
            label: "Dashboard",
            path: "/technician",
        },


        {
            label: "Customer Management",
            path: "/technician/customers",
        },

        {
            label: "Job Management",
            path: "/technician/jobs",
        },

        {
            label: "My Profile",
            path: "/technician/profile",
        },

    ],

    head_technician: [
        
        {
            label: "Dashboard",
            path: "/head-technician",
        },

        {
            label: "Customer Management",
            path: "/head-technician/customers",
        },

        {
            label: "Job Management",
            path: "/head-technician/jobs",
        },

        {
            label: "Staff Management",
            path: "/head-technician/staff",
        },

        {
            label: "My Profile",
            path: "/head-technician/profile",
        },

    ],

    admin: [

         {
            label: "Dashboard",
            path: "/admin",
        },

        {
            label: "Customer Management",
            path: "/admin/customers",
        },

        {
            label: "Job Management",
            path: "/admin/jobs",
        },

        {
            label: "Staff Management",
            path: "/admin/staff",
        },

        {
            label: "My Profile",
            path: "/admin/profile",
        },
        {
            label: "Audit Logs",
            path: "/admin/audit-logs",
        },

    ],

    super_admin: [

         {
            label: "Dashboard",
            path: "/super-admin",
        },

        {
            label: "Customer Management",
            path: "/super-admin/customers",
        },

        {
            label: "Job Management",
            path: "/super-admin/jobs",
        },

        {
            label: "Staff Management",
            path: "/super-admin/staff",
        },

        {
            label: "My Profile",
            path: "/super-admin/profile",
        },
        {
            label: "Audit Logs",
            path: "/super-admin/audit-logs",
        },

    ],

};