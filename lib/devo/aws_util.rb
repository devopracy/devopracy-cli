# wrapper module for the ruby aws sdk
require 'inifile'
require 'aws-sdk'
require 'json'

module Aws_Util
  # set up credentialing for aws
  def self.get_creds(account)
  end

  # use a generalized method to create
  # a resource or a client with the creds
  def self.connect_aws(con_type, account)
    creds = get_creds(account)
    if con_type(client)
    else
    end
  end

  # format tags for filtering resources
  def self.format_tags(tag_hash)
  end

  # resource managers
  def self.get_ami(account, tag_hash)
  end

  def self.get_instance(account, tag_hash)
  end

  def self.get_eip(account, tag_hash)
  end

  def self.get_subnet(account, environment)
  end

  # deploy helpers
  def self.assoc_address(instance_id, eip, account)
  end

  def self.decommission(instance_id, account)
  end

  # inventory and monitor methods
  def self.show_cloud(cloud_name, account)
  end

  def self.inventory(resource_name, account)
  end

